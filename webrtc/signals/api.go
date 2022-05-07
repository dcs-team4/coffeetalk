package signals

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

// Optional configuration for TLS on the signaling server.
type TLSConfig struct {
	TLS      bool   // Whether TLS is enabled.
	CertFile string // Path to TLS certificate file.
	KeyFile  string // Path to TLS key file.
}

// Starts a WebRTC signaling server on the given port, optionally with TLS.
func StartServer(port string) {
	http.HandleFunc("/", connectSocket)

	// Runs the web server (with TLS in if in production) until an error occurs.
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

// HTTP handler for establishing a WebSocket connection to the server.
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

	// Adds the requesting participant to the list of participants.
	userID := addUser(user)
	log.Printf("Connection established with client ID %v.\n", userID)

	// Sets up a channel to stop listening on close.
	stop := make(chan struct{}, 1)

	// Spawns a goroutine for handling messages from the user.
	go user.Listen(stop)

	// Adds a handler for removing the participant when the socket is closed.
	socket.SetCloseHandler(func(code int, text string) error {
		stop <- struct{}{}
		removeUser(userID)
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
	socket.WriteJSON(ConnectionSuccessMessage{
		BaseMessage:      BaseMessage{Type: MsgConnectionSuccess},
		ParticipantCount: participantCount,
	})
}
