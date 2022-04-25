package quiz

import (
	"github.com/dcs-team4/coffeetalk/stm"
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

// Returns a handler for listening to MQTT messages.
// When a start message is sent on the appropriate quiz topic,
// triggers the Start event on the given quiz state machine.
func (machine *QuizMachine) StartQuizHandler() events.OnMessage {
	return func(client events.Client, packet events.Packet) (events.Packet, error) {
		if packet.TopicName == QuizStatusTopic && string(packet.Payload) == QuizStartMessage {
			machine.Start <- stm.Trigger{}
		}

		return packet, nil
	}
}
