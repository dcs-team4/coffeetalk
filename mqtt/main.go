package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcs-team4/coffeetalk/mqtt/broker"
	"github.com/dcs-team4/coffeetalk/mqtt/quiz"
)

// Ports for the MQTT broker.
const (
	socketPort = "1882"
	tcpPort    = "1883"
)

func main() {
	// Sets up channel to keep running the server until a cancel signal is received.
	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM)

	// Starts MQTT broker.
	broker := broker.Start(socketPort, tcpPort)
	fmt.Printf("Running MQTT broker on port %v (WebSocket) and %v (TCP)...\n", socketPort, tcpPort)

	// Creates a new quiz state machine, runs it in a goroutine, and listens for start messages.
	quizmachine := quiz.NewMachine(broker)
	go quizmachine.Run()
	broker.Events.OnMessage = quizmachine.StartQuizHandler()
	fmt.Println("Running quiz machine...")

	// Waits for received cancel signal.
	<-cancelSignal

	broker.Close()
	fmt.Println("Server closed")
}
