package quizmachine

const maxQuestionCount = 5

type Question struct {
	Question string
	Answer   string
}

func (question Question) IsNone() bool {
	return question.Question == ""
}

func newQuestion() Question {
	return Question{
		Question: "What is the meaning of life?",
		Answer:   "42",
	}
}
