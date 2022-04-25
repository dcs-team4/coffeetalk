package broker

import (
	"fmt"
	"log"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
)

// Starts an MQTT broker on the given port in a separate goroutine, and returns the server instance.
func Start(port string) *mqtt.Server {
	// Configures the broker server.
	server := mqtt.NewServer(nil)
	tcp := listeners.NewTCP("tcp1", fmt.Sprintf(":%s", port))
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(auth.Allow),
	})
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
