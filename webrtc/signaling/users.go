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
	ID     int
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

	return user, userID
}

// Adds the given user to the users map, giving them a unique userID and returning it.
func (user *User) Register() (userID int) {
	users.Lock.Lock()
	defer users.Lock.Unlock()

	userID = len(users.Map)
	user.ID = userID
	users.Map[userID] = user
	return userID
}

// Returns whether the given user has joined the peer-to-peer stream.
func (user *User) IsPeer() bool {
	user.Lock.RLock()
	defer user.Lock.RUnlock()

	return user.Name != ""
}

// Returns number of users who have joined the peer-to-peer stream.
func (users *Users) PeerCount() int {
	users.Lock.RLock()
	defer users.Lock.RUnlock()

	count := 0
	for _, user := range users.Map {
		if user.IsPeer() {
			count++
		}
	}

	return count
}

// Joins the peer-to-peer stream with the given username, and notifies other users.
// Returns error if joining failed.
func (user *User) JoinPeers(username string) error {
	err := user.SetName(username)
	if err != nil {
		return err
	}

	for _, otherUser := range users.Map {
		if user.ID == otherUser.ID {
			continue
		}

		otherUser.Socket.WriteJSON(PeerJoinedMessage{Message{MsgPeerJoined}, user.ID, username})
	}

	return nil
}

// Sets the name of the user to the given username. Returns error if it fails.
func (user *User) SetName(username string) error {
	if username == "" {
		return errors.New("username cannot be blank")
	}

	user.Lock.Lock()
	defer user.Lock.Unlock()

	user.Name = username

	return nil
}

// Removes the user from the peer-to-peer stream, and notifies other users.
func (user *User) LeavePeers() {
	user.Lock.Lock()
	user.Name = ""
	user.Lock.Unlock()

	user.HandlePeerLeft()
}

// Sends a message to all other users that the given user has left the peer-to-peer stream.
func (user *User) HandlePeerLeft() {
	users.Lock.RLock()
	defer users.Lock.RUnlock()

	for _, otherUser := range users.Map {
		if user.ID == otherUser.ID {
			continue
		}

		otherUser.Socket.WriteJSON(PeerLeftMessage{Message{MsgPeerLeft}, user.ID})
	}
}

// Removes the user with the given userID from the users map, and notifies other peers.
func (user *User) remove() {
	users.Lock.Lock()
	delete(users.Map, user.ID)
	users.Lock.Unlock()

	if user.IsPeer() {
		user.HandlePeerLeft()
	}
}

// Closes the user's WebSocket connection.
func (user *User) closeConnection() {
	user.Lock.Lock()
	defer user.Lock.Unlock()

	if err := user.Socket.Close(); err != nil {
		log.Printf(
			"Failed to close connection for user '%s' (ID %d): %v\n",
			user.Name,
			user.ID,
			err,
		)
	}
}

// Returns the user of the given userID from the users map, or ok=false if not found.
func (users *Users) Get(userID int) (user *User, ok bool) {
	users.Lock.RLock()
	defer users.Lock.RUnlock()

	user, ok = users.Map[userID]
	return user, ok
}
