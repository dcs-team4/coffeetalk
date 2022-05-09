package signals

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

// Starts a WebRTC signaling server on the given port.
func StartServer(port string) {
	http.HandleFunc("/", connectSocket)

	// Runs the signaling server (with TLS in if in production) until an error occurs.
	var err error
	if os.Getenv("ENV") == "production" {
		err = listenAndServeTLS(":"+port, nil)
	} else {
		err = http.ListenAndServe(":"+port, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}

// HTTP handler for clients establishing a WebSocket connection to the server.
func connectSocket(res http.ResponseWriter, req *http.Request) {
	// Upgrades the request to a WebSocket connection.
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Accepts all origins, in order to enable home clients.
		CheckOrigin: func(*http.Request) bool { return true },
	}
	socket, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		http.Error(res, "unable to establish websocket connection", http.StatusInternalServerError)
		return
	}

	user := &User{
		Name:   "",
		Socket: socket,
		Lock:   new(sync.RWMutex),
	}

	// Register the requesting participant in the list of participants.
	userID := user.Register()
	log.Printf("Connection established with client ID %v.\n", userID)

	// Starts a goroutine for handling messages from the user.
	go user.Listen()

	// Adds a handler for removing the participant when the socket is closed.
	socket.SetCloseHandler(func(code int, text string) error {
		removeUser(userID)

		if user.Name != "" {
			user.HandleUserLeft()
		}

		log.Printf("Socket with client ID %v closed.\n", userID)
		return nil
	})

	// Sends connection success message back, with number of participants in call.
	users.Lock.RLock()
	defer users.Lock.RUnlock()
	participantCount := 0
	for _, user := range users.Map {
		if user.Name != "" {
			participantCount++
		}
	}
	socket.WriteJSON(ConnectionSuccessMessage{BaseMessage{MsgConnectionSuccess}, participantCount})
}
