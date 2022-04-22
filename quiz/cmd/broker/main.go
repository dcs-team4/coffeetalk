package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
)

func main() {
	cancelSignal := listenSignal()
	fmt.Println("Starting server...")

	// Configures server.
	server := mqtt.NewServer(nil)
	tcp := listeners.NewTCP("tcp1", ":1883")
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(auth.Allow),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Logs new connections and messages.
	server.Events.OnConnect = func(client events.Client, packet events.Packet) {
		log.Printf("%v connected!\n", client.ID)
	}
	server.Events.OnMessage = func(client events.Client, packet events.Packet) (
		events.Packet, error,
	) {
		log.Printf("Message received on topic \"%v\"\n", packet.TopicName)
		return packet, nil
	}

	// Starts the server in a new goroutine.
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Started!")

	// Waits for received cancel signal.
	<-cancelSignal

	server.Close()
	fmt.Println("Finished")

}

// Listens for cancelling system signals.
// Sends on the returned channel when a signal is received.
func listenSignal() <-chan struct{} {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cancelSignal := make(chan struct{}, 1)
	go func() {
		<-sigs
		fmt.Println("Caught signal")
		cancelSignal <- struct{}{}
	}()

	return cancelSignal
}
