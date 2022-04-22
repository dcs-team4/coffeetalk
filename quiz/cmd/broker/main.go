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
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	fmt.Println("Mochi MQTT Server initializing... TCP")

	// An example of configuring various server options...
	options := &mqtt.Options{
		BufferSize:      0, // Use default values
		BufferBlockSize: 0, // Use default values
	}

	server := mqtt.NewServer(options)
	tcp := listeners.NewTCP("t1", ":1883")
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(auth.Allow),
	})
	if err != nil {
		log.Fatal(err)
	}

	server.Events.OnConnect = func(client events.Client, packet events.Packet) {
		fmt.Printf("%v connected!\n", client.ID)
	}

	server.Events.OnMessage = func(client events.Client, packet events.Packet) (newPacket events.Packet, err error) {
		fmt.Printf("Message received on topic \"%v\"\n", packet.TopicName)

		return packet, nil
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Started!")

	<-done
	fmt.Println("Caught Signal")

	server.Close()
	fmt.Println("Finished")

}
