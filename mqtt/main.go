package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcs-team4/coffeetalk/mqtt/broker"
	"github.com/dcs-team4/coffeetalk/mqtt/quiz"
	mqtt "github.com/mochi-co/mqtt/server"
)

func main() {
	socketPort, tcpPort := getEnv()

	// Sets up channel to keep running the server until it crashes, or a cancel signal is received.
	close := make(chan struct{}, 1)
	go handleCancelSignal(close)

	mqttBroker := runBroker(socketPort, tcpPort, close)
	runQuizMachine(mqttBroker, close)

	// Waits until server cancels/crashes.
	<-close

	mqttBroker.Close()
	log.Println("Server closed.")
}

// Runs MQTT broker concurrently on the given ports, and returns it.
// Sends on the given close channel if it crashes.
func runBroker(socketPort string, tcpPort string, close chan<- struct{}) *mqtt.Server {
	mqttBroker, err := broker.New(socketPort, tcpPort)
	if err != nil {
		log.Panicln(err)
	}

	err = mqttBroker.Serve()
	if err != nil {
		log.Println("MQTT broker failed:", err)
		close <- struct{}{}
	}

	log.Printf("MQTT broker listening on ports %v (WebSocket), %v (TCP)...\n", socketPort, tcpPort)
	return mqttBroker
}

// Runs quiz state machine concurrently, and listens for quiz start messages on the given broker.
// Sends on the given close channel if it crashes.
func runQuizMachine(mqttBroker *mqtt.Server, close chan<- struct{}) {
	quizmachine := quiz.NewMachine(mqttBroker)

	go func() {
		err := quizmachine.Run()
		log.Println("Quiz state machine failed:", err)
		close <- struct{}{}
	}()

	mqttBroker.Events.OnMessage = quizmachine.StartQuizHandler()
	log.Println("Running quiz state machine...")
}

// Gets ports from environment variables.
func getEnv() (socketPort string, tcpPort string) {
	socketPort = os.Getenv("SOCKET_PORT")
	if socketPort == "" {
		socketPort = "1882"
	}

	tcpPort = os.Getenv("TCP_PORT")
	if tcpPort == "" {
		tcpPort = "1883"
	}

	return socketPort, tcpPort
}

// Sends on the given listener channel when a system cancel signal is received.
func handleCancelSignal(listener chan<- struct{}) {
	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM)

	sig := <-cancelSignal
	log.Println("Received cancel signal:", sig)

	listener <- struct{}{}
}
