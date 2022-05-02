package main

import (
	"os"

	"github.com/dcs-team4/coffeetalk/webrtc/signals"
)

func main() {
	// Gets PORT environment variable, defaulting to 8000 if not present.
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	// Gets ENV environment variable, defaulting to "production" if not present.
	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "production"
	}

	// Configures TLS for the server if in a production environment.
	var tlsConfig signals.TLSConfig
	if env == "production" {
		tlsConfig = signals.TLSConfig{
			TLS:      true,
			CertFile: "tls-cert.pem",
			KeyFile:  "tls-key.pem",
		}
	} else {
		tlsConfig.TLS = false
	}

	// Starts the WebRTC signaling server.
	signals.StartServer(port, tlsConfig)
}
