import { messages, sendWebRTCMessage } from "./socket.js";
import { createPeerVideoElement, decrementParticipantCount, displayLogin, DOM } from "../dom.js";

/**
 * Type for object containing WebRTC connection to peer in the stream, and the video element for
 * that peer's stream, along with the video's container element.
 * Video and container are undefined until initialized.
 * @typedef {Object} Peer
 * @property {RTCPeerConnection} connection
 * @property {HTMLElement} [videoContainer]
 * @property {HTMLVideoElement} [video]
 */

/**
 * State of this user's current peer connections, mapping usernames to peers.
 * @type {{ [username: string]: Peer }}
 */
const peers = {};

/**
 * Tries to get a peer of the given username.
 * @param {string} username
 * @returns {(Peer & { ok: true }) | { ok: false }}
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

  sendWebRTCMessage({
    type: messages.VIDEO_ANSWER,
    to: sender,
    sdp: peer.connection.localDescription,
  });
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

/** Closes every peer connection, removes their video elements, and clears the peers state. */
export function closePeerConnections() {
  for (const peer of Object.values(peers)) {
    peer.connection.close();
    if (peer.videoContainer) {
      DOM.videoContainer().removeChild(peer.videoContainer);
    }
  }

  for (const peerName in peers) {
    delete peers[peerName];
  }
}

/**
 * Closes peer connection with the given peer, removing their video and updating the peers state.
 * @param {string} peerName
 */
export function closePeerConnection(peerName) {
  const peer = getPeer(peerName);
  if (!peer.ok) return;

  peer.connection.close();
  if (peer.videoContainer) DOM.videoContainer().removeChild(peer.videoContainer);
  delete peers[peerName];
}

/**
 * Intitializes a peer connection with the given peer name, and returns the peer.
 * @param {string} peerName
 * @returns {Peer}
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
  connection.addEventListener("icecandidate", (event) => handleICECandidate(event, peerName));
  connection.addEventListener("track", (event) => handleTrack(event, peer, peerName));
  connection.addEventListener("negotiationneeded", () => handleNegotiationNeeded(peer, peerName));
  connection.addEventListener("iceconnectionstatechange", () =>
    handleICEConnectionStateChange(peer)
  );
  connection.addEventListener("icegatheringstatechange", () => handleICEGatheringStateChange(peer));
  connection.addEventListener("signalingstatechange", () => handleSignalingStateChange(peer));

  // Adds the peer to the peers state.
  peers[peerName] = peer;

  return peer;
}

/**
 * Handles incoming candidates for ICE (Interactive Connectivity Establishment protocol).
 * @param {RTCPeerConnectionIceEvent} event, @param {string} peerName
 */
function handleICECandidate(event, peerName) {
  if (event.candidate) {
    sendWebRTCMessage({
      type: messages.ICE_CANDIDATE,
      to: peerName,
      sdp: event.candidate,
    });
  }
}

/**
 * Handles event for performing ICE negotiation.
 * @param {Peer} peer, @param {string} peerName
 */
async function handleNegotiationNeeded(peer, peerName) {
  try {
    const offer = await peer.connection.createOffer();
    await peer.connection.setLocalDescription(offer);
  } catch (error) {
    console.log(error);
    return;
  }

  sendWebRTCMessage({
    type: messages.VIDEO_OFFER,
    to: peerName,
    sdp: peer.connection.localDescription,
  });
}

/**
 * Handles incoming video stream on the peer connection.
 * @param {RTCTrackEvent} event, @param {Peer} peer, @param {string} peerName
 */
function handleTrack(event, peer, peerName) {
  if (!peer.video) {
    const { video, container } = createPeerVideoElement(peerName);
    peer.video = video;
    peer.videoContainer = container;
  }

  peer.video.srcObject = event.streams[0];
}

/** @param {Peer} peer */
function handleICEConnectionStateChange(peer) {}

/** @param {Peer} peer */
function handleICEGatheringStateChange(peer) {}

/** @param {Peer} peer */
function handleSignalingStateChange(peer) {}