package signals

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Forever listens for WebSocket messages from the given user, and forwards them to HandleMessage.
func (user *User) Listen() {
	for {
		_, message, err := user.Socket.ReadMessage()
		if err != nil {
			log.Printf("Could not read message: %v\n", err)
			continue
		}

		user.HandleMessage(message)
	}
}

// Handles the incoming message from the given user.
func (user *User) HandleMessage(rawMessage []byte) {
	// First deserializes the message to a BaseMessage for checking its message type.
	var baseMessage BaseMessage
	err := json.Unmarshal(rawMessage, &baseMessage)
	if err != nil {
		errMsg := "Invalid message"
		user.Socket.WriteJSON(NewErrorMessage(errMsg))
		log.Printf("%v: %v\n", errMsg, err)
		return
	}

	switch baseMessage.Type {
	// Fallthrough because video offer, answer and ICE Candidate messages are all handled the same.
	case MsgVideoOffer:
		fallthrough
	case MsgVideoAnswer:
		fallthrough
	case MsgICECandidate:
		// Deserializes the message to a video exchange message.
		var videoExchange VideoExchangeMessage
		if !DeserializeMsg(rawMessage, &videoExchange, user) {
			return
		}

		// Validates the message.
		target, ok := videoExchange.Validate(user)
		if !ok {
			return
		}

		// Forwards the video offer message to the intended target.
		target.WriteJSON(videoExchange)
	case MsgJoinStream:
		var joinStream NameMessage
		if !DeserializeMsg(rawMessage, &joinStream, user) {
			return
		}

		user.JoinStream(joinStream.Username)
	case MsgLeaveStream:
		user.LeaveStream()
	}
}

// Deserializes the given raw message into the given pointer.
// If deserialization failed, sends an error message to the given sender and returns false.
// Otherwise returns true.
func DeserializeMsg(rawMessage []byte, pointer any, sender *User) (ok bool) {
	err := json.Unmarshal(rawMessage, pointer)

	if err != nil {
		errMsg := "Invalid message received"
		sender.Socket.WriteJSON(NewErrorMessage(errMsg))
		log.Printf("%v: %v\n", errMsg, err)
		return false
	}

	return true
}

// Validates the given TargetedMessage. If valid, sets the message's From field to the given
// sender's username, and returns the target's connection. Otherwise, returns ok=false.
func (msg *TargetedMessage) Validate(sender *User) (target *websocket.Conn, ok bool) {
	targetUser, ok := users.GetByName(msg.To)

	if !ok {
		errMsg := "Invalid message target"
		sender.Socket.WriteJSON(NewErrorMessage(errMsg))
		log.Printf("%v: %v\n", errMsg, msg.To)
		return nil, false
	}

	msg.From = sender.Name

	return targetUser.Socket, true
}
