// Package broker wraps around the mochi-co/mqtt library to set up an MQTT broker.
package broker

import (
	_ "embed"
	"log"
	"os"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
)

// Embedded file with TLS certificate.
//go:embed tls-cert.pem
var tlsCertificate []byte

// Embedded file with TLS key.
//go:embed tls-key.pem
var tlsKey []byte

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

	env := os.Getenv("ENV")

	if env == "production" {
		config.TLS = &listeners.TLS{
			Certificate: tlsCertificate,
			PrivateKey:  tlsKey,
		}
	}

	return config
}
