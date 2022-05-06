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

	env := os.Getenv("ENV")

	// Runs the server until an error is encountered.
	var err error
	if env == "production" {
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

	// Spawns a goroutine for handling messages from the user.
	go user.Listen()

	// Adds a handler for removing the participant when the socket is closed.
	socket.SetCloseHandler(func(code int, text string) error {
		removeUser(userID)
		return nil
	})

	res.Write([]byte("websocket connection established"))
}
