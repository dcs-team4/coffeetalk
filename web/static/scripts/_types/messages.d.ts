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
  type ReceivableMessage = ReceivedPeerExchange | PeerStatus | ConnectionSuccess | Error;

  type PeerExchange = {
    to: string;
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
    from: string;
  };

  type JoinPeers = {
    type: MessageTypes["JOIN_PEERS"];
    username: string;
  };

  type LeavePeers = {
    type: MessageTypes["LEAVE_PEERS"];
  };

  type PeerStatus = {
    type: MessageTypes["PEER_JOINED" | "PEER_LEFT"];
    username: string;
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
