package quiz

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Embedded file with a list of questions.
//go:embed questions.json
var questionsJson []byte

// Global list of questions. Should be populated with questionsJson if empty.
var questions = make([]Question, 0)

// The number of questions that should be asked in a quiz session before ending it.
const maxQuestionCount = 5

// Time to wait before moving between question, answer and next question.
const (
	questionDuration = 30 * time.Second
	answerDuration   = 10 * time.Second
)

// Questions and corresponding answers that make up the quiz.
// Includes an ID to check for question uniqueness, and json tags for reading from file.
type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// Selects a question pseudorandomly from the global questions list, excluding any question in the
// given alreadyAsked list. Returns error if it failed to read questions, or all are already asked.
func newQuestion(alreadyAsked []Question) (Question, error) {
	// If questions are uninitialized, reads from file first.
	if len(questions) == 0 {
		err := readQuestions()
		if err != nil {
			return Question{}, err
		}
	}

	// Filters out already asked questions.
	notAsked := make([]Question, 0)
outerLoop:
	for _, askable := range questions {
		for _, asked := range alreadyAsked {
			if askable.ID == asked.ID {
				continue outerLoop
			}
		}

		notAsked = append(notAsked, askable)
	}

	if len(notAsked) == 0 {
		return Question{}, errors.New("tried to get new question when all were asked")
	}

	// Uses pseudo-random seed to select question.
	rand.Seed(time.Now().UnixNano())
	randomQuestion := notAsked[rand.Intn(len(notAsked))]
	return randomQuestion, nil
}

// Reads questions from the embedded questions file, and sets the global questions variable to it.
// Returns error if file reading failed, or questions are misconfigured.
func readQuestions() error {
	err := json.Unmarshal(questionsJson, &questions)
	if err != nil {
		return fmt.Errorf("failed to read questions.json: %w", err)
	}

	if len(questions) < maxQuestionCount {
		return fmt.Errorf(
			"too few configured questions: have %v, want %v", len(questions), maxQuestionCount,
		)
	}

	return nil
}
