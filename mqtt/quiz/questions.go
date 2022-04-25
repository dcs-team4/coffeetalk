package quiz

import (
	_ "embed"
	"encoding/json"
	"log"
	"math/rand"
)

// Embedded file with a list of questions.
//go:embed questions.json
var questionsJson []byte

// Global list of questions.
// If empty, should be populated with questionsJson.
var questions = make([]Question, 0)

// The number of questions that should be asked in a quiz session before ending it.
const maxQuestionCount = 5

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
}
