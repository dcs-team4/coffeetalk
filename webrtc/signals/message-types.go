package signals

// Message types allowed to and from this server.
const (
	MsgVideoOffer        string = "video-offer"
	MsgVideoAnswer       string = "video-answer"
	MsgICECandidate      string = "new-ice-candidate"
	MsgJoinStream        string = "join-stream"
	MsgConnectionSuccess string = "connection-success"
	MsgUserJoined        string = "user-joined"
	MsgLeaveStream       string = "leave-stream"
	MsgUserLeft          string = "user-left"
	MsgError             string = "error"
)

// Base struct to embed in all message types.
type BaseMessage struct {
	// The kind of message that is sent.
	// Messages that share the same type should always have the same structure.
	Type string `json:"type"`
}

// Utility struct to embed in messages that are to be sent directly from one user to another.
type TargetedMessage struct {
	// Username of sender.
	From string `json:"from"`

	// Username of destination.
	To string `json:"to"`
}

// Message sent between two clients when initiating a video stream between each other.
// The ICE (Interactive Connectivity Establishment) protocol is used to negotiate the connection
// between two peers. This message allows clients to exchange ICE candidates between each other.
type VideoExchangeMessage struct {
	BaseMessage     // Type: MsgVideoOffer/MsgVideoAnswer/MsgICECandidate.
	TargetedMessage // Sender and receiver of video offer/answer.

	// Session Description Protocol string: Describes a client's video stream config in order to
	// properly set up a peer-to-peer stream.
	SDP any `json:"sdp"`
}

// Message with a username field, for when a new user wants to join the stream, to signal to other
// users that a new user has joined, or to signal to other users that a user has left.
type NameMessage struct {
	BaseMessage // Type: MsgJoinStream/MsgUserJoined/MsgUserLeft

	// The username the user wishes to use in the stream.
	Username string `json:"username"`
}

// Message sent when a user successfully establishes a socket connection to the server.
type ConnectionSuccessMessage struct {
	BaseMessage // Type: MsgConnectionSuccess

	ParticipantCount int `json:"participantCount"`
}

// Message sent from server to client when a received message causes an error.
type ErrorMessage struct {
	BaseMessage // Type: MsgError

	ErrorMessage string `json:"errorMessage"`
}

// Utility function for constructing a new ErrorMessage.
func NewErrorMessage(errMsg string) ErrorMessage {
	return ErrorMessage{
		BaseMessage:  BaseMessage{Type: MsgError},
		ErrorMessage: errMsg,
	}
}
