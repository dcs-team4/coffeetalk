package main

import (
	"github.com/dcs-team4/coffeetalk/webrtc/signals"
)

const signalingServerPort = "8000"

func main() {
	// Starts the WebRTC signaling server.
	signals.StartServer(signalingServerPort)
}
