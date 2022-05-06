package main

import (
	"log"
	"os"

	"github.com/dcs-team4/coffeetalk/webrtc/signals"
)

func main() {
	// Gets PORT environment variable, defaulting to 8000 if not present.
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	// Starts the WebRTC signaling server.
	log.Printf("WebRTC signaling server listening on port %v...\n", port)
	signals.StartServer(port)
}
