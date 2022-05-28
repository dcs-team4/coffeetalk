package signaling

import (
	"errors"
	"log"
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
	Name   string
	Socket *websocket.Conn
	Lock   *sync.RWMutex
}

// Creates a new user from the given socket connection.
// Adds the user to the map of users, and adds a handler for removing them on socket close.
// Also starts a goroutine to listen to their messages.
func NewUser(socket *websocket.Conn) (user *User, userID int) {
	user = &User{
		Name:   "",
		Socket: socket,
		Lock:   new(sync.RWMutex),
	}

	// Stores the new connection.
	userID = user.Register()

	// Starts a goroutine for handling messages from the user.
	go user.Listen()

	socket.SetCloseHandler(func(code int, text string) error {
		removeUser(userID)

		if user.InStream() {
			user.HandleStreamLeft()
		}

		log.Printf("Socket with client ID %v closed.\n", userID)
		return nil
	})

	return user, userID
}

// Adds the given user to the users map, giving them a unique userID and returning it.
func (user *User) Register() (userID int) {
	users.Lock.Lock()
	defer users.Lock.Unlock()

	userID = len(users.Map)
	users.Map[userID] = user
	return userID
}

// Returns whether the given user has joined the peer-to-peer stream.
func (user *User) InStream() bool {
	user.Lock.RLock()
	defer user.Lock.RUnlock()

	return user.Name != ""
}

// Returns number of users who have joined the peer-to-peer stream.
func (users *Users) InStream() int {
	users.Lock.RLock()
	defer users.Lock.RUnlock()

	count := 0
	for _, user := range users.Map {
		if user.InStream() {
			count++
		}
	}

	return count
}

// Joins the peer-to-peer stream with the given username, and notifies other users.
// Returns error if joining failed.
func (user *User) JoinStream(username string) error {
	err := user.SetName(username)
	if err != nil {
		return err
	}

	for _, otherUser := range users.Map {
		if otherUser.Name == username {
			continue
		}

		otherUser.Socket.WriteJSON(PeerStatusMessage{Message{MsgPeerJoined}, username})
	}

	return nil
}

// Sets the name of the user to the given username. Returns error if it fails.
func (user *User) SetName(username string) error {
	if username == "" {
		return errors.New("username cannot be blank")
	}

	users.Lock.RLock()
	defer users.Lock.RUnlock()

	for _, user := range users.Map {
		if user.Name == username {
			return errors.New("username already taken")
		}
	}

	user.Lock.Lock()
	defer user.Lock.Unlock()

	user.Name = username

	return nil
}

// Removes the user from the peer-to-peer stream, and notifies other users.
func (user *User) LeaveStream() {
	user.Lock.Lock()
	defer user.Lock.Unlock()

	user.HandleStreamLeft()

	user.Name = ""
}

// Sends a message to all other users that the given user has left the peer-to-peer stream.
func (user *User) HandleStreamLeft() {
	users.Lock.RLock()
	defer users.Lock.RUnlock()

	for _, otherUser := range users.Map {
		if otherUser.Name == user.Name {
			continue
		}

		otherUser.Socket.WriteJSON(PeerStatusMessage{Message{MsgPeerLeft}, user.Name})
	}
}

// Removes the user with the given userID from the users map.
func removeUser(userID int) {
	users.Lock.Lock()
	defer users.Lock.Unlock()

	delete(users.Map, userID)
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
