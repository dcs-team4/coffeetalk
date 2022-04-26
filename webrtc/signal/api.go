package signal

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// HTTP handler for establishing a WebSocket connection to the server.
func connectSocket(res http.ResponseWriter, req *http.Request) {
	// Gets and validates the request body.
	var body struct {
		Username string `json:"username"`
	}
	data, err := io.ReadAll(req.Body)
	err = json.Unmarshal(data, &body)
	if err != nil {
		http.Error(res, "invalid request body", http.StatusBadRequest)
		return
	}

	username := body.Username

	// Ensures unique usernames.
	for existingUsername := range users.Map {
		if username == existingUsername {
			http.Error(res, "username already taken", http.StatusBadRequest)
			return
		}
	}

	// Upgrades the request to a WebSocket connection.
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Accepts all origins for now, in order to enable home clients.
		CheckOrigin: func(*http.Request) bool { return true },
	}
	socket, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		http.Error(res, "unable to establish websocket connection", http.StatusInternalServerError)
		return
	}

	user := &User{
		Name:     username,
		Socket:   socket,
		InStream: false,
		Lock:     new(sync.RWMutex),
	}

	// Adds the requesting participant to the list of participants.
	users.Map[username] = user

	// Spawns a goroutine for handling messages from the user.
	go user.Listen()

	// Adds a handler for removing the participant when the socket is closed.
	socket.SetCloseHandler(func(code int, text string) error {
		users.Lock.Lock()
		defer users.Lock.Unlock()
		delete(users.Map, username)
		return nil
	})

	res.Write([]byte("websocket connection established"))
}
