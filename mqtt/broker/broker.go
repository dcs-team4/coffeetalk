// Package broker wraps around the mochi-co/mqtt library to set up an MQTT broker.
package broker

import (
	"embed"
	"log"
	"os"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
)

// For production environment: expects tls-cert.pem and tls-key.pem in tls directory.
//go:embed all:tls
var tlsFiles embed.FS

// Starts an MQTT broker that listens for WebSocket and TCP connections on the given ports.
func Start(socketPort string, tcpPort string) *mqtt.Server {
	// Configures the broker server.
	server := mqtt.NewServer(nil)

	listenerConfig := configureListener()

	// Listens for WebSocket connections on the given socketPort.
	socket := listeners.NewWebsocket("socket1", ":"+socketPort)
	err := server.AddListener(socket, listenerConfig)

	// Listens for TCP connections on the given tcpPort.
	tcp := listeners.NewTCP("tcp1", ":"+tcpPort)
	err = server.AddListener(tcp, listenerConfig)

	if err != nil {
		log.Fatal(err)
	}

	// Logs new connections.
	server.Events.OnConnect = func(client events.Client, packet events.Packet) {
		log.Printf("%v connected!\n", client.ID)
	}

	// Starts the server in a new goroutine.
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return server
}

// Returns an MQTT listener config.
// If in a production environment, configures TLS for the listener with the embedded TLS files.
func configureListener() *listeners.Config {
	config := &listeners.Config{Auth: new(auth.Allow)}

	if os.Getenv("ENV") == "production" {
		tlsCertificate, err := tlsFiles.ReadFile("tls/tls-cert.pem")
		if err != nil {
			log.Printf("TLS certificate setup failed: %v\n", err)
			return config
		}

		tlsKey, err := tlsFiles.ReadFile("tls/tls-key.pem")
		if err != nil {
			log.Printf("TLS key setup failed: %v\n", err)
			return config
		}

		config.TLS = &listeners.TLS{
			Certificate: tlsCertificate,
			PrivateKey:  tlsKey,
		}
	}

	return config
}
