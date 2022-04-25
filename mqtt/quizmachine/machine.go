package quizmachine

import (
	"time"

	"github.com/dcs-team4/coffeetalk/mqtt/messages"
	"github.com/dcs-team4/coffeetalk/stm"
	mqtt "github.com/mochi-co/mqtt/server"
)

// State machine for quiz sessions.
type QuizMachine struct {
	// Triggered to start a new quiz session.
	Start stm.Event

	// Triggered when the time between question and answer has run out.
	questionTimer stm.Event

	// Triggered when the time between an answer and the next question has run out.
	answerTimer stm.Event

	// List of questions asked so far in the current quiz session.
	questions []Question

	// The MQTT broker where questions and answers are to be published.
	broker *mqtt.Server
}

// IDs of the quiz machine's states.
// Uses iota for automatic enumeration (0, 1, 2...).
const (
	idleState stm.StateID = iota
	questionState
	answerState
)

// Map of quiz machine state IDs to the functions that should run for those states.
var states = stm.States[*QuizMachine]{
	idleState:     IdleState,
	questionState: QuestionState,
}

// Returns a new quiz state machine, with all channels and lists initialized.
// Attaches the given broker to the machine, and assumes it is valid to send on.
func New(broker *mqtt.Server) *QuizMachine {
	return &QuizMachine{
		Start:         make(stm.Event),
		questionTimer: make(stm.Event),
		answerTimer:   make(stm.Event),
		questions:     make([]Question, 0),
		broker:        broker,
	}
}

// Runs the given quiz state machine. Keeps running through every configured state function,
// transitioning to new states as they return.
func (machine *QuizMachine) Run() {
	currentState := idleState

	for {
		currentState = states[currentState](machine)
	}
}

// Gets the latest question to process in the quiz session.
// Assumes that there is at least one question in the questions list.
func (machine QuizMachine) currentQuestion() Question {
	if len(machine.questions) == 0 {
		return Question{}
	}

	return machine.questions[len(machine.questions)-1]
}

// Waits for a Start event to trigger, then returns the Question state as the next state.
func IdleState(machine *QuizMachine) (nextState stm.StateID) {
	<-machine.Start
	return questionState
}

// Adds a new question to the machine's questions list, publishes it to the MQTT broker,
// then waits 30 seconds before returning the Answer state as the next state.
func QuestionState(machine *QuizMachine) (nextState stm.StateID) {
	machine.questions = append(machine.questions, newQuestion())

	machine.broker.Publish(messages.QuestionTopic, []byte(machine.currentQuestion().Question), true)

	go stm.SetTimer(30*time.Second, machine.questionTimer)
	<-machine.questionTimer
	return answerState
}

// Publishes the answer to the previous question to the MQTT broker. Then, if the quiz has reached
// its final question, ends the quiz and returns the Idle state; otherwise, waits 10 seconds and
// then returns the Question state as the next state.
func AnswerState(machine *QuizMachine) (nextState stm.StateID) {
	machine.broker.Publish(messages.AnswerTopic, []byte(machine.currentQuestion().Answer), true)

	// If the quiz has reached its final question, sends the quiz end message,
	// resets the quiz questions, and returns to the idle state.
	if len(machine.questions) >= maxQuestionCount {
		machine.broker.Publish(messages.QuizStatusTopic, []byte(messages.QuizEndMessage), true)
		machine.questions = make([]Question, 0)
		return idleState
	}

	go stm.SetTimer(10*time.Second, machine.answerTimer)
	<-machine.answerTimer
	return questionState
}
