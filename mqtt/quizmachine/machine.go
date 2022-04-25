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

func New(broker *mqtt.Server) *QuizMachine {
	return &QuizMachine{
		Start:         make(stm.Event),
		questionTimer: make(stm.Event),
		answerTimer:   make(stm.Event),
		questions:     make([]Question, 0),
		broker:        broker,
	}
}

func (machine *QuizMachine) Run() {
	currentState := idleState

	for {
		currentState = states[currentState](machine)
	}
}

func (machine QuizMachine) currentQuestion() Question {
	if len(machine.questions) == 0 {
		return Question{}
	}

	return machine.questions[len(machine.questions)-1]
}

func IdleState(machine *QuizMachine) (nextState stm.StateID) {
	// Waits for a start event, then returns the Question state as the next state.
	<-machine.Start
	return questionState
}

func QuestionState(machine *QuizMachine) (nextState stm.StateID) {
	machine.questions = append(machine.questions, newQuestion())

	// Publishes the current question to the MQTT broker.
	machine.broker.Publish(messages.QuestionTopic, []byte(machine.currentQuestion().Question), true)

	// Waits 30 seconds before showing the answer to the question.
	go stm.SetTimer(30*time.Second, machine.questionTimer)
	<-machine.questionTimer
	return answerState
}

func AnswerState(machine *QuizMachine) (nextState stm.StateID) {
	// Publishes the answer to the current question on the MQTT broker.
	machine.broker.Publish(messages.AnswerTopic, []byte(machine.currentQuestion().Answer), true)

	// If the quiz has reached its final question, sends the quiz end message,
	// resets the quiz questions, and returns to the idle state.
	if len(machine.questions) >= maxQuestionCount {
		machine.broker.Publish(messages.QuizStatusTopic, []byte(messages.QuizEndMessage), true)
		machine.questions = make([]Question, 0)
		return idleState
	}

	// Waits 10 seconds, then moves on to the next question.
	go stm.SetTimer(10*time.Second, machine.answerTimer)
	<-machine.answerTimer
	return questionState
}
