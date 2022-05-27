/** Type declarations for messages as defined by the WebRTC signaling server. */
declare namespace messages {
  type MessageTypes = {
    VIDEO_OFFER: "video-offer";
    VIDEO_ANSWER: "video-answer";
    ICE_CANDIDATE: "ice-candidate";
    JOIN_STREAM: "join-stream";
    LEAVE_STREAM: "leave-stream";
    PEER_JOINED: "peer-joined";
    PEER_LEFT: "peer-left";
    CONNECTION_SUCCESS: "connection-success";
    ERROR: "error";
  };

  /** Messages that the WebRTC server expects the client to send. */
  type SendableMessage = PeerExchange | JoinStream | LeaveStream;

  /** Messages that the WebRTC server expects the client to receive. */
  type ReceivableMessage = ReceivedPeerExchange | PeerStatus | ConnectionSuccess | Error;

  type PeerExchange = {
    to: string;
  } & (
    | {
        type: MessageTypes["VIDEO_OFFER" | "VIDEO_ANSWER"];
        data: RTCSessionDescription;
      }
    | {
        type: MessageTypes["ICE_CANDIDATE"];
        data: RTCIceCandidateInit;
      }
  );

  type ReceivedPeerExchange = PeerExchange & { from: string };

  type JoinStream = {
    type: MessageTypes["JOIN_STREAM"];
    username: string;
  };

  type LeaveStream = {
    type: MessageTypes["LEAVE_STREAM"];
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
