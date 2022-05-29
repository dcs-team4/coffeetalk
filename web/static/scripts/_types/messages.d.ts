/** Type declarations for messages as defined by the WebRTC signaling server. */
declare namespace messages {
  type MessageTypes = {
    PEER_OFFER: "peer-offer";
    PEER_ANSWER: "peer-answer";
    ICE_CANDIDATE: "ice-candidate";
    JOIN_PEERS: "join-peers";
    LEAVE_PEERS: "leave-peers";
    PEER_JOINED: "peer-joined";
    PEER_LEFT: "peer-left";
    CONNECTION_SUCCESS: "connection-success";
    ERROR: "error";
  };

  /** Messages that the WebRTC server expects the client to send. */
  type SendableMessage = PeerExchange | JoinPeers | LeavePeers;

  /** Messages that the WebRTC server expects the client to receive. */
  type ReceivableMessage = ReceivedPeerExchange | PeerJoined | PeerLeft | ConnectionSuccess | Error;

  type PeerExchange = {
    receiverId: number;
  } & (
    | {
        type: MessageTypes["PEER_OFFER" | "PEER_ANSWER"];
        data: RTCSessionDescriptionInit;
      }
    | {
        type: MessageTypes["ICE_CANDIDATE"];
        data: RTCIceCandidateInit;
      }
  );

  type ReceivedPeerExchange = PeerExchange & {
    senderId: number;
    senderName: string;
  };

  type JoinPeers = {
    type: MessageTypes["JOIN_PEERS"];
    name: string;
  };

  type LeavePeers = {
    type: MessageTypes["LEAVE_PEERS"];
  };

  type PeerJoined = {
    type: MessageTypes["PEER_JOINED"];
    id: number;
    name: string;
  };

  type PeerLeft = {
    type: MessageTypes["PEER_LEFT"];
    id: number;
  };

  type ConnectionSuccess = {
    type: MessageTypes["CONNECTION_SUCCESS"];
    peerCount: number;
  };

  type Error = {
    type: MessageTypes["ERROR"];
    errorMessage: string;
  };
}
