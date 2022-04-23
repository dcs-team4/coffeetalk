package quizmachine

import (
	"time"

	"github.com/dcs-team4/coffeetalk/mqtt/messages"
	"github.com/dcs-team4/coffeetalk/stm"
	mqtt "github.com/mochi-co/mqtt/server"
)

type QuizMachine struct {
	Start stm.Event

	post stm.Event

	questionTimer stm.Event

	answerTimer stm.Event

	questionCount int

	currentQuestion Question

	broker *mqtt.Server
}

// IDs of the quiz machine's states.
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
		post:          make(stm.Event),
		questionCount: 0,
		broker:        broker,
	}
}

func (machine *QuizMachine) Run() {
	currentState := idleState

	for {
		currentState = states[currentState](machine)
	}
}

func IdleState(machine *QuizMachine) (nextState stm.StateID) {
	// Waits for a start event, then returns the Question state as the next state.
	<-machine.Start
	return questionState
}

func QuestionState(machine *QuizMachine) (nextState stm.StateID) {
	machine.currentQuestion = getQuestion()

	// Publishes the current question to the MQTT broker.
	machine.broker.Publish(messages.QuestionTopic, []byte(machine.currentQuestion.Question), true)

	// Waits 30 seconds before showing the answer to the question.
	go stm.SetTimer(30*time.Second, machine.questionTimer)
	<-machine.questionTimer
	return answerState
}

func AnswerState(machine *QuizMachine) (nextState stm.StateID) {
	// Publishes the answer to the current question on the MQTT broker.
	machine.broker.Publish(messages.AnswerTopic, []byte(machine.currentQuestion.Answer), true)

	machine.questionCount++

	// If the quiz has reached its final question, sends the quiz end message,
	// resets the quiz variables, and returns to the idle state.
	if machine.questionCount >= maxQuestionCount {
		machine.broker.Publish(messages.QuizStatusTopic, []byte(messages.QuizEndMessage), true)
		machine.currentQuestion = Question{}
		machine.questionCount = 0
		return idleState
	}

	// Waits 10 seconds, then moves on to the next question.
	go stm.SetTimer(10*time.Second, machine.answerTimer)
	<-machine.answerTimer
	return questionState
}

func (machine QuizMachine) Listen() {
	started := false
	questionCount := 0

	for {
		select {
		case <-machine.Start:
			if started {
				break
			}

			started = true
			machine.post <- stm.Trigger{}
		case <-machine.post:
			if !started {
				break
			}

			if questionCount == maxQuestionCount {
				started = false
				break
			}

			questionCount++
		}
	}
}
