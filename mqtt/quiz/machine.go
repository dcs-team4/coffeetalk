package quiz

import (
	"errors"
	"fmt"

	"github.com/dcs-team4/coffeetalk/stm"
	mqtt "github.com/mochi-co/mqtt/server"
)

// State machine for quiz sessions.
// Implements stm.StateMachine.
type QuizMachine struct {
	// Map of quiz machine state IDs to the functions that should run for those states.
	states stm.States[*QuizMachine]

	// Triggered to start a new quiz session.
	start stm.Event

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
// Uses iota for automatic enumeration, +1 to avoid potential clash with zero value.
const (
	idleState stm.StateID = iota + 1
	questionState
	answerState
)

// Returns a new quiz state machine, with all states, channels and lists initialized.
// Attaches the given broker to the machine, and assumes it is valid to send on.
func NewMachine(broker *mqtt.Server) *QuizMachine {
	return &QuizMachine{
		states: stm.States[*QuizMachine]{
			idleState:     runIdleState,
			questionState: runQuestionState,
			answerState:   runAnswerState,
		},
		start:         make(stm.Event),
		questionTimer: make(stm.Event),
		answerTimer:   make(stm.Event),
		questions:     make([]Question, 0),
		broker:        broker,
	}
}

// Returns the quiz machine's configured states.
func (machine *QuizMachine) States() stm.States[*QuizMachine] {
	return machine.states
}

// Runs the given quiz state machine. Keeps running through every configured state function,
// transitioning to new states as they return, until an error occurs.
func (machine *QuizMachine) Run() error {
	startState := idleState
	err := stm.RunMachine(machine, startState)
	return err
}

// Gets the latest question to process in the quiz session. Returns error if question list is empty.
func (machine QuizMachine) currentQuestion() (Question, error) {
	if len(machine.questions) == 0 {
		return Question{}, errors.New("tried to get question from empty question list")
	}

	return machine.questions[len(machine.questions)-1], nil
}

// Waits for a Start event to trigger, then returns the Question state as the next state.
func runIdleState(machine *QuizMachine) (nextState stm.StateID, err error) {
	<-machine.start
	return questionState, nil
}

// Adds a new question to the machine's questions list, publishes it to the MQTT broker,
// then waits for the question timer event before transitioning to the Answer state.
func runQuestionState(machine *QuizMachine) (nextState stm.StateID, err error) {
	question, err := newQuestion(machine.questions)
	if err != nil {
		return 0, fmt.Errorf("failed to get new quiz question: %w", err)
	}
	machine.questions = append(machine.questions, question)

	machine.broker.Publish(QuestionTopic, []byte(question.Question), true)

	go stm.SetTimer(questionDuration, machine.questionTimer)
	<-machine.questionTimer

	return answerState, nil
}

// Publishes the answer to the previous question to the MQTT broker.
// Then, if the quiz has reached its final question, ends the quiz and returns the Idle state;
// otherwise, waits for the answer timer event before transitioning to the Question state.
func runAnswerState(machine *QuizMachine) (nextState stm.StateID, err error) {
	question, err := machine.currentQuestion()
	if err != nil {
		return 0, fmt.Errorf("quiz machine answer state failed: %w", err)
	}

	machine.broker.Publish(AnswerTopic, []byte(question.Answer), true)

	// If this is the final question, end the quiz and clean up the questions.
	if len(machine.questions) >= maxQuestionCount {
		machine.broker.Publish(QuizStatusTopic, []byte(QuizEndMessage), true)
		machine.questions = make([]Question, 0)
		return idleState, nil
	}

	go stm.SetTimer(answerDuration, machine.answerTimer)
	<-machine.answerTimer

	return questionState, nil
}
