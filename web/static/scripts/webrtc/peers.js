import { sendWebRTCMessage, messages } from "./socket.js";
import { DOM, createPeerVideoElement } from "../dom.js";

/**
 * State of this user's current peer connections, mapping usernames to peers.
 * @type {{ [username: string]: types.Peer }}
 */
const peers = {};

/**
 * Tries to get a peer of the given username.
 * @param {string} username
 * @returns {(types.Peer & { ok: true }) | { ok: false }}
 */
function getPeer(username) {
  if (!peers.hasOwnProperty(username)) {
    console.log(`Unrecognized peer: ${username}`);
    return { ok: false };
  }

  return { ...peers[username], ok: true };
}

/**
 * This user's video stream. Undefined until initialized.
 * @type {MediaStream | undefined}
 */
let localStream;

/** Tries to get the user's video and audio stream, and if successful, displays it.  */
export async function streamLocalVideo() {
  try {
    localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
  } catch (error) {
    console.log("Failed to get local video stream:", error);
    return;
  }

  DOM.localVideo().srcObject = localStream;
}

/**
 * Attempts to stream this user's video to the given peer connection.
 * @param {RTCPeerConnection} peerConnection
 */
function sendLocalStream(peerConnection) {
  if (!localStream) {
    console.log("Failed to send local stream: Stream uninitialized.");
    return;
  }

  for (const track of localStream.getTracks()) {
    peerConnection.addTrack(track, localStream);
  }
}

/**
 * Creates a peer connection with a peer of the given name,
 * and starts sending this user's stream to them.
 * @param {string} peerName
 */
export async function sendVideoOffer(peerName) {
  const peer = createPeerConnection(peerName);
  sendLocalStream(peer.connection);
}

/**
 * Handles a received video stream offer from a peer.
 * @param {string} sender Peer username.
 * @param {RTCSessionDescriptionInit} sdp Session Description Protocol data.
 */
export async function receiveVideoOffer(sender, sdp) {
  const peer = createPeerConnection(sender);

  // Configures the WebRTC session.
  const session = new RTCSessionDescription(sdp);
  await peer.connection.setRemoteDescription(session);

  sendLocalStream(peer.connection);

  // Attempts to create an answer to the connection offer.
  try {
    const answer = await peer.connection.createAnswer();
    await peer.connection.setLocalDescription(answer);
  } catch (error) {
    console.log("Error in establishing WebRTC connection:", error);
    return;
  }

  if (peer.connection.localDescription) {
    sendWebRTCMessage({
      type: messages.VIDEO_ANSWER,
      to: sender,
      sdp: peer.connection.localDescription,
    });
  }
}

/**
 * Handles a received answer to a video stream offer, finalizing the WebRTC connection.
 * @param {string} sender Peer username.
 * @param {RTCSessionDescriptionInit} sdp Session Description Protocol data.
 */
export async function receiveVideoAnswer(sender, sdp) {
  const peer = getPeer(sender);
  if (!peer.ok) return;

  // Finalizes peer connection configuration.
  const session = new RTCSessionDescription(sdp);
  await peer.connection.setRemoteDescription(session);
}

/**
 * Handles received ICE (Interactive Connectivity Establishment protocol) candidate for setting up
 * WebRTC connection.
 * @param {string} sender Peer username.
 * @param {RTCIceCandidateInit} sdp SDP candidate used for ICE.
 */
export async function receiveICECandidate(sender, sdp) {
  const peer = getPeer(sender);
  if (!peer.ok) return;

  const candidate = new RTCIceCandidate(sdp);

  try {
    await peer.connection.addIceCandidate(candidate);
  } catch (error) {
    console.log(error);
  }
}

/**
 * Closes peer connection with the given peer, removing their video and updating the peers state.
 * @param {string} peerName
 */
export function closePeerConnection(peerName) {
  const peer = getPeer(peerName);
  if (!peer.ok) return;

  // Stops the video stream.
  if (peer.video) {
    peer.video.pause();
    for (const track of /** @type {MediaStream} */ (peer.video.srcObject)?.getTracks()) {
      track.stop();
    }
  }

  // Remove video from the DOM.
  if (peer.videoContainer) {
    DOM.videos().removeChild(peer.videoContainer);
  }

  // Stop peer-to-peer communication.
  for (const transceiver of peer.connection.getTransceivers()) {
    transceiver.stop();
  }

  // Close connection and clear peer.
  peer.connection.close();
  delete peers[peerName];
}

/** Closes every peer connection. */
export function closePeerConnections() {
  for (const peerName in peers) {
    closePeerConnection(peerName);
  }
}

/**
 * Intitializes a peer connection with the given peer name, and returns the peer.
 * @param {string} peerName
 * @returns {types.Peer}
 */
export function createPeerConnection(peerName) {
  // Initializes the peer connection with a STUN service (Session Traversal Utilities for NAT).
  const connection = new RTCPeerConnection({
    iceServers: [
      {
        // Open source STUN service (https://www.stunprotocol.org/)
        urls: "stun:stun.stunprotocol.org",
      },
    ],
  });
  const peer = { connection };

  // Registers event handlers.
  connection.addEventListener("track", (event) => {
    handleTrack(event, peer, peerName);
  });
  connection.addEventListener("negotiationneeded", () => {
    handleNegotiationNeeded(peer, peerName);
  });
  connection.addEventListener("icecandidate", (event) => {
    handleICECandidate(event, peerName);
  });
  connection.addEventListener("iceconnectionstatechange", () => {
    handleICEConnectionStateChange(peer, peerName);
  });
  connection.addEventListener("signalingstatechange", () => {
    handleSignalingStateChange(peer, peerName);
  });

  // Adds the peer to the peers state.
  peers[peerName] = peer;

  return peer;
}

/**
 * Handles incoming video stream on the peer connection.
 * @param {RTCTrackEvent} event, @param {types.Peer} peer, @param {string} peerName
 */
function handleTrack(event, peer, peerName) {
  if (!peer.video) {
    const { video, container } = createPeerVideoElement(peerName);
    peer.video = video;
    peer.videoContainer = container;
  }

  peer.video.srcObject = event.streams[0];
}

/**
 * Handles event for performing ICE negotiation.
 * @param {types.Peer} peer, @param {string} peerName
 */
async function handleNegotiationNeeded(peer, peerName) {
  try {
    const offer = await peer.connection.createOffer();
    await peer.connection.setLocalDescription(offer);
  } catch (error) {
    console.log("Error in ICE negotiation:", error);
    return;
  }

  if (peer.connection.localDescription) {
    sendWebRTCMessage({
      type: messages.VIDEO_OFFER,
      to: peerName,
      sdp: peer.connection.localDescription,
    });
  }
}

/**
 * Handles incoming candidates for ICE (Interactive Connectivity Establishment protocol).
 * @param {RTCPeerConnectionIceEvent} event, @param {string} peerName
 */
function handleICECandidate(event, peerName) {
  if (event.candidate) {
    sendWebRTCMessage({ type: messages.ICE_CANDIDATE, to: peerName, sdp: event.candidate });
  }
}

/**
 * Handles changes in the ICE connection state with the given peer.
 * @param {types.Peer} peer, @param {string} peerName
 */
function handleICEConnectionStateChange(peer, peerName) {
  switch (peer.connection.iceConnectionState) {
    case "closed":
    case "failed":
      closePeerConnection(peerName);
      break;
  }
}

/**
 * Handles changes in the signaling state with the given peer.
 * @param {types.Peer} peer, @param {string} peerName
 */
function handleSignalingStateChange(peer, peerName) {
  if (peer.connection.signalingState === "closed") {
    closePeerConnection(peerName);
  }
}
