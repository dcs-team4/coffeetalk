package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcs-team4/coffeetalk/mqtt/broker"
	"github.com/dcs-team4/coffeetalk/mqtt/quiz"
	"hermannm.dev/ipfinder"
)

func main() {
	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM)

	socketPort := "8080"
	tcpPort := "1883"
	broker := broker.Start(socketPort, tcpPort)
	fmt.Printf("Running MQTT broker (TCP port: %v, WebSocket port: %v)\n", tcpPort, socketPort)

	localIPs, err := ipfinder.FindLocalIPs()
	if err == nil {
		for network, ips := range localIPs {
			for _, ip := range ips {
				fmt.Printf("IP: (%v) %v\n", network, ip)
			}
		}
	}

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
