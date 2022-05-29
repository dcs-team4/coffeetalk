// Package stm provides utility types and functions for setting up state machines in Go.
package stm

import (
	"fmt"
	"time"
)

// A state machine is a type with a configured set of states and corresponding state functions,
// and a method for running the machine.
// The type parameter here is to avoid self-reference, but typically points back to the implementer.
type StateMachine[Machine any] interface {
	// Returns the configured states for this state machine.
	States() States[Machine]

	// Manages the machine's states, and executes its state functions until an error occurs.
	// Typically calls the RunMachine function provided by this package.
	Run() error
}

// Uniquely identifies a state in a state machine.
// Typically used with iota for easy enumeration.
type StateID int

// A function to run when in a given state.
// Returns when transitioning states (typically after an event triggers), or when an error occurs.
// Takes a type parameter for the type of state machine to execute on.
type StateFunc[Machine any] func(Machine) (nextState StateID, err error)

// A map of possible states for a state machine.
// Takes a type parameter for the type of state machine to attach to.
type States[Machine any] map[StateID]StateFunc[Machine]

// A signal that may trigger a change in a state machine, typically caused by outside input.
// Implemented as an empty struct in order to take 0 space.
type Trigger struct{}

// An event is something that can happen when a state machine receives a trigger.
// Implemented as a channel, to allow for concurrent communication with the state machine,
// and listening for the event.
type Event chan Trigger

// Utility function for running a state machine.
// Keeps running every configured state function (starting with the given startState),
// and transitions to new states as they return. Keeps running until an error occurs.
// The type parameter constraint ensures that the machine configures states for itself.
func RunMachine[Machine StateMachine[Machine]](machine Machine, startState StateID) error {
	currentState := startState

	for {
		// Gets the state function for the current state from the machine's state configuration.
		currentStateFunc, ok := machine.States()[currentState]
		if !ok {
			return fmt.Errorf("missing state machine function for state ID %v", currentState)
		}

		// Runs the state function, and sets its returned next state as the new current state.
		nextState, err := currentStateFunc(machine)
		if err != nil {
			return fmt.Errorf(
				"error in state machine function for state ID %v: %w", currentState, err,
			)
		}
		currentState = nextState
	}
}

// Utility function for setting a timer and triggering the given state machine event on expiry.
// Typically run in a goroutine, to later listen on the event.
func SetTimer(duration time.Duration, event Event) {
	time.Sleep(duration)
	event <- Trigger{}
}
