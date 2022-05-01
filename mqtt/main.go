package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/dcs-team4/coffeetalk/mqtt/broker"
	"github.com/dcs-team4/coffeetalk/mqtt/quiz"
	"github.com/joho/godotenv"
)

func main() {
	// Loads environment variables from .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load .env file.")
	}

	// Get ports from environment variables.
	socketPort := os.Getenv("PORT")
	if socketPort == "" {
		socketPort = "1882"
	}
	tcpPort := os.Getenv("TCP_PORT")
	if tcpPort == "" {
		tcpPort = "1883"
	}

	// Sets up channel to keep running the server until a cancel signal is received.
	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt, os.Kill)

	// Starts MQTT broker.
	broker := broker.Start(socketPort, tcpPort)
	log.Printf("Running MQTT broker on port %v (WebSocket) and %v (TCP)...\n", socketPort, tcpPort)

	// Creates a new quiz state machine, runs it in a goroutine, and listens for start messages.
	quizmachine := quiz.NewMachine(broker)
	go quizmachine.Run()
	broker.Events.OnMessage = quizmachine.StartQuizHandler()
	log.Println("Running quiz machine...")

	// Waits for received cancel signal.
	sig := <-cancelSignal

	log.Println(sig.String())
	broker.Close()
	log.Println("Server closed.")
}
