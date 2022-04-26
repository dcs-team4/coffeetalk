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
	case MsgVideoOffer:
		// Falls through to the next case, as video offer and answer messages are handled the same.
		fallthrough
	case MsgVideoAnswer:
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
	case MsgICECandidate:
		// Deserializes the message to an ICE candidate message.
		var iceCandidate ICECandidateMessage
		if !DeserializeMsg(rawMessage, &iceCandidate, user) {
			return
		}

		// Validates the message.
		target, ok := iceCandidate.Validate(user)
		if !ok {
			return
		}

		// Forwards the ICE candidate message to the intended target.
		target.WriteJSON(iceCandidate)
	case MsgJoinStream:
		user.Lock.Lock()
		user.InStream = true
		user.Lock.Unlock()
	case MsgLeaveStream:
		user.Lock.Lock()
		user.InStream = false
		user.Lock.Unlock()
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
	targetUser, ok := users.Get(msg.To)

	if !ok {
		errMsg := "Invalid message target"
		sender.Socket.WriteJSON(NewErrorMessage(errMsg))
		log.Printf("%v: %v\n", errMsg, msg.To)
		return nil, false
	}

	msg.From = sender.Name

	return targetUser.Socket, true
}