package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcs-team4/coffeetalk/quiz/broker"
)

func main() {
	cancelSignal := listenSignal()
	fmt.Println("Starting server...")

	broker := broker.Start("1883")

	// Waits for received cancel signal.
	<-cancelSignal

	broker.Close()
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
