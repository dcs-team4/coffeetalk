package signaling

import (
	"encoding/json"
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
		var peerExchange PeerExchangeMessage
		if !DeserializeMsg(rawMessage, &peerExchange, user) {
			return
		}

		target, ok := peerExchange.Validate(user)
		if !ok {
			return
		}

		// Forwards the peer exchange message to the intended target.
		target.WriteJSON(peerExchange)
	case MsgJoinPeers:
		var joinStream PeerStatusMessage
		if !DeserializeMsg(rawMessage, &joinStream, user) {
			return
		}

		log.Printf("%v message received from: %v\n", MsgJoinPeers, joinStream.Username)

		user.JoinPeers(joinStream.Username)
	case MsgLeavePeers:
		log.Printf("%v message received from: %v\n", MsgLeavePeers, user.Name)

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

// Validates the given peer exchange message. If valid, sets the message's From field to the given
// sender's username, and returns the target peer's connection. Otherwise, returns ok=false.
func (message *PeerExchangeMessage) Validate(sender *User) (target *websocket.Conn, ok bool) {
	targetUser, ok := users.GetByName(message.To)

	if !ok {
		errMsg := "Invalid message target"
		sender.Socket.WriteJSON(ErrorMessage{Message{MsgError}, errMsg})
		log.Printf("%v: %v\n", errMsg, message.To)
		return nil, false
	}

	sender.Lock.RLock()
	defer sender.Lock.RUnlock()
	message.From = sender.Name

	return targetUser.Socket, true
}
