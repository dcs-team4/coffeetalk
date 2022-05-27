package quiz

import (
	"log"

	"github.com/dcs-team4/coffeetalk/stm"
	"github.com/mochi-co/mqtt/server/events"
)

// Quiz-related constants for MQTT communication.
const (
	// Name of the MQTT topic where the server posts quiz questions for clients to see.
	QuestionTopic string = "coffeetalk/quiz/questions"

	// Name of the MQTT topic where the server posts answers to quiz questions.
	AnswerTopic string = "coffeetalk/quiz/answers"

	// Name of the MQTT topic where clients post to start a quiz,
	// and the server posts if the quiz is ended.
	QuizStatusTopic string = "coffeetalk/quiz/status"

	// The message posted on the MQTT quiz status topic to start a quiz.
	QuizStartMessage string = "start-quiz"

	// The message posted on the MQTT quiz status topic when a quiz ends.
	QuizEndMessage string = "end-quiz"
)

// Returns a handler for listening to MQTT messages.
// When a start message is sent on the appropriate quiz topic,
// triggers the Start event on the given quiz state machine.
func (machine *QuizMachine) StartQuizHandler() events.OnMessage {
	return func(client events.Client, packet events.Packet) (events.Packet, error) {
		log.Printf(
			"Message received (topic: %v, message: %v)\n",
			packet.TopicName,
			string(packet.Payload),
		)

		if packet.TopicName == QuizStatusTopic && string(packet.Payload) == QuizStartMessage {
			machine.Start <- stm.Trigger{}
		}

		return packet, nil
	}
}
