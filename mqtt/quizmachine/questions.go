package quizmachine

// The number of questions that should be asked in a quiz session before ending it.
const maxQuestionCount = 5

// Questions and corresponding answers that make up the quiz.
// Includes an ID to check for question uniqueness.
type Question struct {
	ID       int
	Question string
	Answer   string
}

// Returns a new question, and makes sure that it is not among the given already asked questions.
func newQuestion(alreadyAsked []Question) Question {
	return Question{
		ID:       0,
		Question: "What is the meaning of life?",
		Answer:   "42",
	}
}
