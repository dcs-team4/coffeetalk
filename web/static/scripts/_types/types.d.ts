/** Utility type declarations. */
declare namespace types {
  type Peer = {
    name: string;
    connection: RTCPeerConnection;
    videoContainer?: HTMLElement;
    video?: HTMLVideoElement;
  };
}
