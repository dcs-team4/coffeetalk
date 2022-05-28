package main

import (
	"log"
	"os"

	"github.com/dcs-team4/coffeetalk/webrtc/signaling"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	log.Printf("WebRTC signaling server listening on port %v...\n", port)

	err := signaling.StartServer(port)
	if err != nil {
		log.Println("Signaling server failed:", err)
	}

	log.Println("WebRTC signaling server closed.")
}
