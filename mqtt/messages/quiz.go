package messages

import (
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
)

// Quiz-related constants for MQTT communication.
const (
	// Name of the MQTT topic where clients post to start a quiz,
	// and the server posts if the quiz is ended.
	QuizStatusTopic string = "TTM4115/t4/quiz/s"

	// The message posted on the MQTT quiz status topic to start a quiz.
	QuizStartMessage string = "Quiz_started"

	// The message posted on the MQTT quiz status topic when a quiz ends.
	QuizEndMessage string = "Quiz_ended"

	// Name of the MQTT topic where the server posts quiz questions for clients to see.
	QuestionTopic string = "TTM4115/t4/quiz/q"

	// Name of the MQTT topic where the server posts answers to quiz questions.
	AnswerTopic string = "TTM4115/t4/quiz/a"
)

// Entry point of the quizmaster.
// Listens to quiz-related messages on the given MQTT broker.
func Listen(broker *mqtt.Server) {
	broker.Events.OnMessage = quizStartHandler(broker)
}

// Returns a handler for listening to messages on the given MQTT broker.
// Starts a quiz if the appropriate message is posted, otherwise does nothing.
func quizStartHandler(broker *mqtt.Server) events.OnMessage {
	return func(client events.Client, packet events.Packet) (events.Packet, error) {
		if packet.TopicName == QuizStatusTopic {
			startQuiz(broker)
		}

		return packet, nil
	}
}

func startQuiz(broker *mqtt.Server) {

}
