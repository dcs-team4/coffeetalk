package signaling

// Message types allowed to and from this server.
const (
	MsgPeerOffer         string = "peer-offer"
	MsgPeerAnswer        string = "peer-answer"
	MsgICECandidate      string = "ice-candidate"
	MsgJoinPeers         string = "join-peers"
	MsgLeavePeers        string = "leave-peers"
	MsgPeerJoined        string = "peer-joined"
	MsgPeerLeft          string = "peer-left"
	MsgConnectionSuccess string = "connection-success"
	MsgError             string = "error"
)

// Base struct to embed in all message types.
type Message struct {
	Type string `json:"type"`
}

// Message sent between two clients when initiating a peer-to-peer connection between each other.
type PeerExchangeMessage struct {
	Message           // Type: MsgPeerOffer/MsgPeerAnswer/MsgICECandidate
	ReceiverID int    `json:"receiverId"`
	SenderID   int    `json:"senderId"`
	SenderName string `json:"senderName"`

	// WebRTC data for setting up peer-to-peer connection. Format depends on message type.
	//
	// For peer offer/answer messages:
	// SDP (Session Description Protocol) object with initial config proposal for the connection.
	//
	// For ICE (Interactive Connectivity Establishment) candidate messages:
	// Candidate object used to negotiate peer-to-peer connection setup.
	//
	// Typed as any, as the server only forwards the message and does not care about its type.
	Data any `json:"data"`
}

// Message sent by client when they want to join the peer-to-peer stream.
type JoinPeersMessage struct {
	Message        // Type: MsgJoinPeers
	Name    string `json:"name"` // The name that the client wants to use in the stream.
}

// Message sent by client when they want to leave the peer-to-peer stream.
// Identical to base Message; included here for documentation purposes.
type LeavePeersMessage struct {
	Message // Type: MsgLeavePeers
}

// Message sent by server to notify users that a new peer has joined the stream.
type PeerJoinedMessage struct {
	Message        // Type: MsgPeerJoined/MsgPeerLeft
	ID      int    `json:"id"`
	Name    string `json:"name"`
}

// Message sent by the server to notify users of a peer leaving.
type PeerLeftMessage struct {
	Message     // Type: MsgPeerLeft
	ID      int `json:"id"`
}

// Message sent when a user successfully establishes a socket connection to the server.
type ConnectionSuccessMessage struct {
	Message       // Type: MsgConnectionSuccess
	PeerCount int `json:"peerCount"` // Number of active peers in the stream.
}

// Message sent from server to client when a received message causes an error.
type ErrorMessage struct {
	Message             // Type: MsgError
	ErrorMessage string `json:"errorMessage"`
}
