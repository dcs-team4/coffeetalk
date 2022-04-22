package quizmaster

import (
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
)

// Quiz-related constants for MQTT communication.
const (
	// Name of the MQTT topic where clients post to start a quiz.
	quizStartTopic string = "TTM4115/t4/quiz/s"

	// The message posted on the MQTT topic to start a quiz.
	quizStartMessage string = "Quiz_started"

	// Name of the MQTT topic where the server posts quiz questions for clients to see.
	questionTopic string = "TTM4115/t4/quiz/q"

	// Name of the MQTT topic where the server posts answers to quiz questions.
	answerTopic string = "TTM4115/t4/quiz/a"
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
		if packet.TopicName == quizStartTopic {
			startQuiz(broker)
		}

		return packet, nil
	}
}

func startQuiz(broker *mqtt.Server) {

}
