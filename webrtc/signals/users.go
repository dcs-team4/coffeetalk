package signals

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

// Global map of users connected to the server.
var users = Users{
	Map:  make(map[int]*User),
	Lock: &sync.RWMutex{},
}

// Map of user IDs to users connected to the server, with a mutex for thread-safe modification.
type Users struct {
	Map  map[int]*User
	Lock *sync.RWMutex
}

// A user connected to the signaling server.
type User struct {
	ID     int
	Name   string
	Socket *websocket.Conn
	Lock   *sync.RWMutex
}

// Returns whether the given user has joined the video stream.
func (user User) InStream() bool {
	return user.Name != ""
}

// Sets the user's name to the given username, and notifies other users that a new user has joined.
func (user *User) JoinStream(username string) error {
	if username == "" {
		return errors.New("Username cannot be blank.")
	}

	users.Lock.RLock()
	defer users.Lock.RUnlock()

	for _, user := range users.Map {
		if user.Name == username {
			return errors.New("Username already taken.")
		}
	}

	user.Lock.Lock()
	defer user.Lock.Unlock()

	user.Name = username

	for _, otherUser := range users.Map {
		if otherUser.Name == username {
			continue
		}

		otherUser.Socket.WriteJSON(PeerStatusMessage{Message{MsgPeerJoined}, username})
	}

	return nil
}

// Removes the given user from the stream, and signals to other users that the user has left.
func (user *User) LeaveStream() {
	user.Lock.Lock()
	defer user.Lock.Unlock()

	user.HandleUserLeft()

	user.Name = ""
}

// Sends a message to all other users that the given user has left the stream.
func (user *User) HandleUserLeft() {
	users.Lock.RLock()
	defer users.Lock.RUnlock()

	for _, otherUser := range users.Map {
		if otherUser.Name == user.Name {
			continue
		}

		otherUser.Socket.WriteJSON(PeerStatusMessage{Message{MsgPeerLeft}, user.Name})
	}
}

// Returns the user of the given userID from the users map, or ok=false if not found.
func (users *Users) GetByID(userID int) (user *User, ok bool) {
	users.Lock.RLock()
	defer users.Lock.RUnlock()

	user, ok = users.Map[userID]
	return user, ok
}

// Returns the user of the given username from the users map, or ok=false if not found.
func (users *Users) GetByName(username string) (user *User, ok bool) {
	users.Lock.RLock()
	defer users.Lock.RUnlock()

	for _, user := range users.Map {
		if user.Name == username {
			return user, true
		}
	}

	return nil, false
}

// Adds the given user to the users map, giving it a unique userID and returning it.
func (user *User) Register() (userID int) {
	users.Lock.Lock()
	defer users.Lock.Unlock()

	userID = len(users.Map)
	user.ID = userID
	users.Map[userID] = user
	return userID
}

// Removes the user with the given userID from the users map.
func removeUser(userID int) {
	users.Lock.Lock()
	defer users.Lock.Unlock()

	delete(users.Map, userID)
}
