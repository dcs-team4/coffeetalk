package quiz

import (
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
	answerState:   AnswerState,
}

// Returns a new quiz state machine, with all channels and lists initialized.
// Attaches the given broker to the machine, and assumes it is valid to send on.
func NewMachine(broker *mqtt.Server) *QuizMachine {
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

// Adds a new question to the machine's questions list, publishes it to the MQTT broker, then waits
// for the question timer event before returning the Answer state as the next state.
func QuestionState(machine *QuizMachine) (nextState stm.StateID) {
	// Adds a new question to the questions list.
	machine.questions = append(machine.questions, newQuestion(machine.questions))

	// Publishes the current question to the MQTT broker.
	machine.broker.Publish(QuestionTopic, []byte(machine.currentQuestion().Question), true)

	// Sets a timer to trigger the question timer event.
	go stm.SetTimer(questionDuration, machine.questionTimer)

	// Waits for the question timer event.
	<-machine.questionTimer

	return answerState
}

// Publishes the answer to the previous question to the MQTT broker. Then, if the quiz has reached
// its final question, ends the quiz and returns the Idle state; otherwise, waits for the answer
// timer event before returning the Question state as the next state.
func AnswerState(machine *QuizMachine) (nextState stm.StateID) {
	// Publishes the answer to the current question to the MQTT broker.
	machine.broker.Publish(AnswerTopic, []byte(machine.currentQuestion().Answer), true)

	// If the quiz has reached its final question, sends the quiz end message,
	// resets the quiz questions, and returns to the idle state.
	if len(machine.questions) >= maxQuestionCount {
		machine.broker.Publish(QuizStatusTopic, []byte(QuizEndMessage), true)
		machine.questions = make([]Question, 0)
		return idleState
	}

	// Sets a timer to trigger the answer timer event.
	go stm.SetTimer(answerDuration, machine.answerTimer)

	// Waits for the answer timer event.
	<-machine.answerTimer

	return questionState
}
