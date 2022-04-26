package signal

import (
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	Name     string
	Socket   *websocket.Conn
	InStream bool
	Lock     *sync.RWMutex
}

type Users struct {
	Map  map[string]*User
	Lock *sync.RWMutex
}

var users = Users{
	Map:  make(map[string]*User),
	Lock: &sync.RWMutex{},
}

func (users *Users) Get(username string) (*User, bool) {
	users.Lock.RLock()
	defer users.Lock.RUnlock()
	user, ok := users.Map[username]
	return user, ok
}
