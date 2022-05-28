package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcs-team4/coffeetalk/mqtt/broker"
	"github.com/dcs-team4/coffeetalk/mqtt/quiz"
)

func main() {
	// Gets ports from environment variables.
	socketPort := os.Getenv("SOCKET_PORT")
	if socketPort == "" {
		socketPort = "1882"
	}
	tcpPort := os.Getenv("TCP_PORT")
	if tcpPort == "" {
		tcpPort = "1883"
	}

	// Sets up channel to keep running the server until a cancel signal is received.
	done := make(chan struct{}, 1)
	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-cancelSignal
		log.Println("Received cancel signal:", sig)
		done <- struct{}{}
	}()

	// Runs MQTT broker concurrently.
	mqttBroker := broker.New(socketPort, tcpPort)
	go func() {
		err := mqttBroker.Serve()
		log.Println("MQTT broker failed:", err)
		done <- struct{}{}
	}()
	log.Printf("MQTT broker listening on ports %v (WebSocket), %v (TCP)...\n", socketPort, tcpPort)

	// Runs quiz state machine concurrently, and listens for quiz start messages on the broker.
	quizmachine := quiz.NewMachine(mqttBroker)
	go func() {
		err := quizmachine.Run()
		log.Println("Quiz state machine failed:", err)
		done <- struct{}{}
	}()
	mqttBroker.Events.OnMessage = quizmachine.StartQuizHandler()
	log.Println("Running quiz state machine...")

	// Waits until server cancels/crashes.
	<-done

	mqttBroker.Close()
	log.Println("Server closed.")
}
