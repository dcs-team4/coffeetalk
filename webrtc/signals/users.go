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

// Map user IDs to users connected to the server, with a mutex for thread-safe modification.
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

	return nil
}

// Removes the given user from the stream.
func (user *User) LeaveStream() {
	user.Lock.Lock()
	defer user.Lock.Unlock()

	user.Name = ""
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
func addUser(user *User) (userID int) {
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
