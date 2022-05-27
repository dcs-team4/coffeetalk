/** Message types as defined by the WebRTC signaling server. */
declare namespace messages {
  type MessageTypes = {
    VIDEO_OFFER: "video-offer";
    VIDEO_ANSWER: "video-answer";
    ICE_CANDIDATE: "new-ice-candidate";
    JOIN_STREAM: "join-stream";
    LEAVE_STREAM: "leave-stream";
    USER_JOINED: "user-joined";
    USER_LEFT: "user-left";
    CONNECTION_SUCCESS: "connection-success";
    ERROR: "error";
  };

  /** Messages that the WebRTC server expects the client to send. */
  type SendableMessage = VideoExchange | ICECandidate | JoinStream | LeaveStream;

  /** Messages that the WebRTC server expects the client to receive. */
  type ReceivableMessage =
    | ReceivedVideoExchange
    | ReceivedICECandidate
    | PeerChange
    | ConnectionSuccess
    | Error;

  type VideoExchange = {
    type: MessageTypes["VIDEO_OFFER" | "VIDEO_ANSWER"];
    to: string;
    sdp: RTCSessionDescriptionInit;
  };

  type ReceivedVideoExchange = VideoExchange & { from: string };

  type ICECandidate = {
    type: MessageTypes["ICE_CANDIDATE"];
    to: string;
    sdp: RTCIceCandidateInit;
  };

  type ReceivedICECandidate = ICECandidate & { from: string };

  type JoinStream = {
    type: MessageTypes["JOIN_STREAM"];
    username: string;
  };

  type LeaveStream = {
    type: MessageTypes["LEAVE_STREAM"];
  };

  type PeerChange = {
    type: MessageTypes["USER_JOINED" | "USER_LEFT"];
    username: string;
  };

  type ConnectionSuccess = {
    type: MessageTypes["CONNECTION_SUCCESS"];
    participantCount: number;
  };

  type Error = {
    type: MessageTypes["ERROR"];
    errorMessage: string;
  };
}
