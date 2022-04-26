// Package broker wraps around the mochi-co/mqtt library to set up an MQTT broker.
package broker

import (
	"fmt"
	"log"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
)

// Starts an MQTT broker that listens for WebSocket and TCP connections
func Start(socketPort string, tcpPort string) *mqtt.Server {
	// Configures the broker server.
	server := mqtt.NewServer(nil)

	// Listens for WebSocket connections on the given socketPort.
	socket := listeners.NewWebsocket("socket1", fmt.Sprintf(":%s", socketPort))
	err := server.AddListener(socket, &listeners.Config{Auth: new(auth.Allow)})

	// Listens for TCP connections on the given tcpPort.
	tcp := listeners.NewTCP("tcp1", fmt.Sprintf(":%s", tcpPort))
	err = server.AddListener(tcp, &listeners.Config{Auth: new(auth.Allow)})

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
	fmt.Println("Broker listening...")

	return server
}
