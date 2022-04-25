package quiz

import (
	_ "embed"
	"encoding/json"
	"log"
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
// given alreadyAsked list.
func newQuestion(alreadyAsked []Question) Question {
	// If questions are uninitialized, reads from file first.
	if len(questions) == 0 {
		readQuestions()
	}

	// Filters out already asked questions.
	notAsked := make([]Question, 0)
outerLoop:
	for _, question := range questions {
		for _, asked := range alreadyAsked {
			if question.ID == asked.ID {
				continue outerLoop
			}
		}

		notAsked = append(notAsked, question)
	}

	// Returns a question at a pseudo-random index.
	return notAsked[rand.Intn(len(notAsked))]
}

// Reads questions from the embedded questions.json file,
// setting the global questions variable to it.
func readQuestions() {
	err := json.Unmarshal(questionsJson, &questions)
	if err != nil {
		log.Fatal("failed to read questions.json")
	}

	// Crashes if too few questions are configured.
	if len(questions) < maxQuestionCount {
		log.Fatalf(
			"too few configured questions; have %v, want %v\n",
			len(questions),
			maxQuestionCount,
		)
	}
}
