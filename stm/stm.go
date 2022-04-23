// Package stm provides utility types for creating state machines in Go.
package stm

import "time"

// A state machine can be essentially any type, up to the implementer.
type Machine interface {
	// Executes the state machine.
	// Typically runs a StateFunc given by the States map.
	Run()
}

// Uniquely identifies a state in a state machine.
// Typically used with iota constants for easy enumeration.
type StateID int

// A function to run to represent a state.
// Takes a type parameter for the type of state machine to execute on.
// Typically keeps running until some event is triggered, upon which it returns the next state.
type StateFunc[M Machine] func(machine M) (nextState StateID)

// A map of possible states for a state machine.
// Takes a type parameter for the type of state machine to attach to.
type States[M Machine] map[StateID]StateFunc[M]

// A signal that may trigger a change in a state machine, typically caused by outside input.
// Implemented as an empty struct in order to take 0 space.
type Trigger struct{}

// An event is something that can happen when a state machine receives a trigger.
// Implemented as a channel, to allow for concurrent communication with the state machine,
// and listening for the event.
type Event chan Trigger

// Utility function for setting a timer and triggering the given state machine event on expiry.
// Typically run in a goroutine, to later listen on the event.
func SetTimer(t time.Duration, event Event) {
	time.Sleep(t)
	event <- Trigger{}
}
