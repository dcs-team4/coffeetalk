// Package broker wraps around the mochi-co/mqtt library to set up an MQTT broker.
package broker

import (
	"embed"
	"fmt"
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

// Creates an MQTT broker configured to listen for WebSocket and TCP connections on the given ports.
// Returns an error if setup failed.
func New(socketPort string, tcpPort string) (*mqtt.Server, error) {
	broker := mqtt.NewServer(nil)

	listenerConfig, err := configureListener()
	if err != nil {
		return nil, err
	}

	socket := listeners.NewWebsocket("socket1", ":"+socketPort)
	err = broker.AddListener(socket, listenerConfig)
	if err != nil {
		return nil, fmt.Errorf("websocket listener setup failed: %w", err)
	}

	tcp := listeners.NewTCP("tcp1", ":"+tcpPort)
	err = broker.AddListener(tcp, listenerConfig)
	if err != nil {
		return nil, fmt.Errorf("tcp listener setup failed: %w", err)
	}

	broker.Events.OnConnect = func(client events.Client, packet events.Packet) {
		log.Printf("%v connected.\n", client.ID)
	}

	return broker, nil
}

// Returns an MQTT listener config, or an error if config setup failed.
// If in a production environment, configures TLS for the listener with the embedded TLS files.
func configureListener() (*listeners.Config, error) {
	config := &listeners.Config{Auth: new(auth.Allow)}

	if os.Getenv("ENV") == "production" {
		tlsCertificate, err := tlsFiles.ReadFile("tls/tls-cert.pem")
		if err != nil {
			return nil, fmt.Errorf("tls certificate setup failed: %w", err)
		}

		tlsKey, err := tlsFiles.ReadFile("tls/tls-key.pem")
		if err != nil {
			return nil, fmt.Errorf("tls key setup failed: %w", err)
		}

		config.TLS = &listeners.TLS{
			Certificate: tlsCertificate,
			PrivateKey:  tlsKey,
		}
	}

	return config, nil
}
