package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcs-team4/coffeetalk/mqtt/broker"
	"github.com/dcs-team4/coffeetalk/mqtt/quiz"
)

func main() {
	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM)

	port := "1883"
	broker := broker.Start(port)
	fmt.Printf("Running MQTT broker on port %v...\n", port)

	// Creates a new quiz machine, runs it in a goroutine, and listens for start messages.
	quizmachine := quiz.NewMachine(broker)
	go quizmachine.Run()
	broker.Events.OnMessage = quizmachine.StartQuizHandler()
	fmt.Println("Running quiz machine...")

	// Waits for received cancel signal.
	<-cancelSignal

	broker.Close()
	fmt.Println("Server closed")
}
