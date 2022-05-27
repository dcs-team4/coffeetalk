package signals

// Message types allowed to and from this server.
const (
	MsgVideoOffer        string = "video-offer"
	MsgVideoAnswer       string = "video-answer"
	MsgICECandidate      string = "ice-candidate"
	MsgJoinStream        string = "join-stream"
	MsgLeaveStream       string = "leave-stream"
	MsgPeerJoined        string = "peer-joined"
	MsgPeerLeft          string = "peer-left"
	MsgConnectionSuccess string = "connection-success"
	MsgError             string = "error"
)

// Base struct to embed in all message types.
type Message struct {
	// The kind of message that is sent.
	// Messages that share the same type should always have the same structure.
	Type string `json:"type"`
}

// Message sent between two clients when initiating a peer-to-peer video stream between each other.
type PeerExchangeMessage struct {
	Message // Type: MsgVideoOffer/MsgVideoAnswer/MsgICECandidate

	// Username of sender.
	From string `json:"from"`

	// Username of receiving peer.
	To string `json:"to"`

	// WebRTC data for setting up peer-to-peer connection. Format depends on message type.
	//
	// For video offer/answer messages:
	// SDP (Session Description Protocol) object with proposed configuration for the call.
	//
	// For ICE (Interactive Connectivity Establishment) candidate messages:
	// Candidate object used to negotiate peer-to-peer connection setup.
	//
	// Typed as any, as the server only forwards the message and does not care about its type.
	Data any `json:"data"`
}

// Message for signaling changes in a user's peer status: when a user wants to join the peer-to-peer
// stream, or to signal to other peers that a peer has joined or left.
type PeerStatusMessage struct {
	Message // Type: MsgJoinStream/MsgPeerJoined/MsgPeerLeft

	// Username of the peer affected by the message.
	Username string `json:"username"`
}

// Message sent when a user successfully establishes a socket connection to the server.
type ConnectionSuccessMessage struct {
	Message // Type: MsgConnectionSuccess

	// Number of active peers in the stream.
	PeerCount int `json:"peerCount"`
}

// Message sent from server to client when a received message causes an error.
type ErrorMessage struct {
	Message // Type: MsgError

	ErrorMessage string `json:"errorMessage"`
}
