package signaling

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Listens for WebSocket messages from the given user, and forwards them to HandleMessage.
// Stops when the socket is closed.
func (user *User) Listen() {
	for {
		_, message, err := user.Socket.ReadMessage()

		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				return
			}
			log.Println("Failed to read socket message:", err)
			continue
		}

		user.HandleMessage(message)
	}
}

// Handles the incoming message from the given user.
func (user *User) HandleMessage(rawMessage []byte) {
	// First deserializes the message to a BaseMessage for checking its message type.
	var baseMessage Message
	err := json.Unmarshal(rawMessage, &baseMessage)
	if err != nil {
		errMsg := "Invalid message"
		user.Socket.WriteJSON(ErrorMessage{Message{MsgError}, errMsg})
		log.Printf("%v: %v\n", errMsg, err)
		return
	}

	// Handles the message according to its type.
	switch baseMessage.Type {
	// Fallthrough because peer offer, answer and ICE candidate messages are all handled the same.
	case MsgPeerOffer:
		fallthrough
	case MsgPeerAnswer:
		fallthrough
	case MsgICECandidate:
		var message PeerExchangeMessage
		if !DeserializeMsg(rawMessage, &message, user) {
			return
		}

		target, ok := message.Validate(user)
		if !ok {
			return
		}

		// Forwards the peer exchange message to the intended target.
		target.WriteJSON(message)
	case MsgJoinPeers:
		var message JoinPeersMessage
		if !DeserializeMsg(rawMessage, &message, user) {
			return
		}

		log.Printf(
			"User %v requested to join the stream with username: %v", user.ID, message.Name,
		)

		err := user.JoinPeers(message.Name)
		if err != nil {
			log.Printf("User %v (%v) failed to join stream: %v\n", user.ID, message.Name, err)
			user.Socket.WriteJSON(ErrorMessage{
				Message{MsgError},
				fmt.Sprint("Failed to join peer-to-peer stream:", err),
			})
		}
	case MsgLeavePeers:
		log.Printf(
			"User %v (%v) requested to leave the stream.\n", user.ID, user.Name,
		)

		user.LeavePeers()
	default:
		log.Printf("Unrecognized message type received: %v\n", baseMessage.Type)
	}
}

// Deserializes the given raw message into the given pointer.
// If deserialization failed, sends an error message to the given sender and returns false.
// Otherwise returns true.
func DeserializeMsg(rawMessage []byte, pointer any, sender *User) (ok bool) {
	err := json.Unmarshal(rawMessage, pointer)

	if err != nil {
		errMsg := "Invalid message received"
		sender.Socket.WriteJSON(ErrorMessage{Message{MsgError}, errMsg})
		log.Printf("%v: %v\n", errMsg, err)
		return false
	}

	return true
}

// Validates the given peer exchange message. If valid, sets the message's sender fields to the
// given sender's ID and name, and returns the receiving peer's connection.
func (message *PeerExchangeMessage) Validate(sender *User) (receiver *websocket.Conn, valid bool) {
	receivingUser, valid := users.Get(message.ReceiverID)

	if !valid {
		errMsg := "Invalid message target"
		sender.Socket.WriteJSON(ErrorMessage{Message{MsgError}, errMsg})
		log.Printf("%v: %v\n", errMsg, message.ReceiverID)
		return nil, false
	}

	sender.Lock.RLock()
	defer sender.Lock.RUnlock()

	message.SenderID = sender.ID
	message.SenderName = sender.Name

	return receivingUser.Socket, true
}
