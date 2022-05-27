/** Utility type declarations. */
declare namespace types {
  type Peer = {
    connection: RTCPeerConnection;
    videoContainer?: HTMLElement;
    video?: HTMLVideoElement;
  };
}
