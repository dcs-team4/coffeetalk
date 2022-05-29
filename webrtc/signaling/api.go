package signaling

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// Starts a WebRTC signaling server on the given port. Keeps running until an error occurs.
func StartServer(port string) error {
	http.HandleFunc("/", connectSocket)

	if os.Getenv("ENV") == "production" {
		return listenAndServeTLS(":"+port, nil)
	} else {
		return http.ListenAndServe(":"+port, nil)
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

	_, userID := NewUser(socket)
	log.Printf("Connection established with user ID %v.\n", userID)

	// Sends connection success message back, with number of peers in call.
	socket.WriteJSON(ConnectionSuccessMessage{Message{MsgConnectionSuccess}, users.PeerCount()})
}
