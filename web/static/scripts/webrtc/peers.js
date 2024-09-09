import { sendWebRTCMessage, messages } from "./socket.js";
import { DOM, createPeerVideoElement } from "../dom.js";

/**
 * State of this user's current peer connections, mapping user IDs to peers.
 * @type {Map<number, types.Peer>}
 */
const peers = new Map();

/**
 * Tries to get a peer of the given username.
 * @param {number} id
 * @returns {(types.Peer & { ok: true }) | { ok: false }}
 */
function getPeer(id) {
  const peer = peers.get(id);

  if (!peer) {
    console.log("Unrecognized peer ID:", id);
    return { ok: false };
  }

  return { ...peer, ok: true };
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
 * Creates a peer connection with a peer of the given ID and name,
 * and starts sending this user's stream to them.
 * @param {number} peerID, @param {string} peerName
 */
export async function handleNewPeer(peerID, peerName) {
  const peer = createPeerConnection(peerID, peerName);
  sendLocalStream(peer.connection);
}

/**
 * Handles a received WebRTC connection offer from a peer.
 * @param {number} senderId
 * @param {string} senderName
 * @param {RTCSessionDescriptionInit} sessionDescription
 */
export async function receivePeerOffer(senderId, senderName, sessionDescription) {
  const peer = createPeerConnection(senderId, senderName);

  // Configures the WebRTC session.
  const session = new RTCSessionDescription(sessionDescription);
  await peer.connection.setRemoteDescription(session);

  sendLocalStream(peer.connection);

  // Attempts to create an answer to the connection offer.
  try {
    const answer = await peer.connection.createAnswer();
    await peer.connection.setLocalDescription(answer);
  } catch (error) {
    console.log("Failed to establish WebRTC connection:", error);
    return;
  }

  if (peer.connection.localDescription) {
    sendWebRTCMessage({
      type: messages.PEER_ANSWER,
      receiverId: senderId,
      data: peer.connection.localDescription,
    });
  }
}

/**
 * Handles a received answer from a peer to a previously sent WebRTC connection offer.
 * @param {number} senderId, @param {RTCSessionDescriptionInit} sessionDescription
 */
export async function receivePeerAnswer(senderId, sessionDescription) {
  const peer = getPeer(senderId);
  if (!peer.ok) return;

  // Finalizes peer connection configuration.
  const session = new RTCSessionDescription(sessionDescription);
  await peer.connection.setRemoteDescription(session);
}

/**
 * Handles ICE (Interactive Connectivity Establishment) candidate for setting up WebRTC connection.
 * @param {number} senderId, @param {RTCIceCandidateInit} iceCandidate
 */
export async function receiveICECandidate(senderId, iceCandidate) {
  const peer = getPeer(senderId);
  if (!peer.ok) return;

  const candidate = new RTCIceCandidate(iceCandidate);

  try {
    await peer.connection.addIceCandidate(candidate);
  } catch (error) {
    console.log("Failed to add received ICE candidate:", error);
  }
}

/**
 * Closes peer connection with the given peer, removing their video and updating the peers state.
 * @param {number} peerId
 */
export function closePeerConnection(peerId) {
  const peer = getPeer(peerId);
  if (!peer.ok) return;

  if (peer.video) {
    peer.video.pause();

    const videoSrc = /** @type {MediaStream | null} */ (peer.video.srcObject);
    for (const track of videoSrc?.getTracks() ?? []) {
      track.stop();
    }
  }

  if (peer.videoContainer) {
    DOM.videos().removeChild(peer.videoContainer);
  }

  for (const transceiver of peer.connection.getTransceivers()) {
    transceiver.stop();
  }

  peer.connection.close();
  peers.delete(peerId);
}

/** Closes every peer connection. */
export function closePeerConnections() {
  for (const peerId of peers.keys()) {
    closePeerConnection(peerId);
  }
}

/**
 * Intitializes a peer connection with the given peer name, and returns the peer.
 * @param {number} peerId, @param {string} peerName
 * @returns {types.Peer}
 */
export function createPeerConnection(peerId, peerName) {
  // Initializes the peer connection with a STUN service (Session Traversal Utilities for NAT).
  const connection = new RTCPeerConnection({
    iceServers: [
      // Using Google's public STUN servers
      {
        urls: "stun:stun.l.google.com:19302",
      },
      {
        urls: "stun:stun1.l.google.com:19302,",
      },
      {
        urls: "stun:stun2.l.google.com:19302",
      },
      {
        urls: "stun:stun3.l.google.com:19302",
      },
      {
        urls: "stun:stun4.l.google.com:19302",
      },
    ],
  });
  const peer = { name: peerName, connection };
  peers.set(peerId, peer);

  connection.addEventListener("track", (event) => {
    handleTrack(event, peer, peerName);
  });
  connection.addEventListener("negotiationneeded", () => {
    handleNegotiationNeeded(peer, peerId);
  });
  connection.addEventListener("icecandidate", (event) => {
    handleICECandidate(event, peerId);
  });
  connection.addEventListener("iceconnectionstatechange", () => {
    handleICEConnectionStateChange(peer, peerId);
  });
  connection.addEventListener("signalingstatechange", () => {
    handleSignalingStateChange(peer, peerId);
  });

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
 * @param {types.Peer} peer, @param {number} peerId
 */
async function handleNegotiationNeeded(peer, peerId) {
  try {
    const offer = await peer.connection.createOffer();
    await peer.connection.setLocalDescription(offer);
  } catch (error) {
    console.log("Failed to create peer connection offer:", error);
    return;
  }

  if (peer.connection.localDescription) {
    sendWebRTCMessage({
      type: messages.PEER_OFFER,
      receiverId: peerId,
      data: peer.connection.localDescription,
    });
  }
}

/**
 * Handles incoming candidates for ICE (Interactive Connectivity Establishment protocol).
 * @param {RTCPeerConnectionIceEvent} event, @param {number} peerId
 */
function handleICECandidate(event, peerId) {
  if (event.candidate) {
    sendWebRTCMessage({ type: messages.ICE_CANDIDATE, receiverId: peerId, data: event.candidate });
  }
}

/**
 * Handles changes in the ICE connection state with the given peer.
 * @param {types.Peer} peer, @param {number} peerId
 */
function handleICEConnectionStateChange(peer, peerId) {
  switch (peer.connection.iceConnectionState) {
    case "closed":
    case "failed":
      closePeerConnection(peerId);
      break;
  }
}

/**
 * Handles changes in the signaling state with the given peer.
 * @param {types.Peer} peer, @param {number} peerId
 */
function handleSignalingStateChange(peer, peerId) {
  if (peer.connection.signalingState === "closed") {
    closePeerConnection(peerId);
  }
}
