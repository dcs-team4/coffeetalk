package quizmachine

// The number of questions that should be asked in a quiz session before ending it.
const maxQuestionCount = 5

// Questions and corresponding answers that make up the quiz.
type Question struct {
	Question string
	Answer   string
}

// Returns a new question.
func newQuestion() Question {
	return Question{
		Question: "What is the meaning of life?",
		Answer:   "42",
	}
}
