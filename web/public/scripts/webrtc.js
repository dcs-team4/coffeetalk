//@ts-check

import { messages, sendToServer } from "./socket.js";

/**
 * @typedef {Object} Peer
 * @property {RTCPeerConnection} connection
 * @property {HTMLVideoElement} [video]
 */

/** @type {{ [key: string]: Peer }} */
const peers = {};

const mediaConstraints = {
  audio: true,
  video: true,
};

export async function sendVideoOffer(peerName) {
  const peer = createPeerConnection(peerName);

  let stream;
  try {
    stream = await navigator.mediaDevices.getUserMedia(mediaConstraints);
  } catch (error) {
    console.log(error);
    return;
  }

  /** @type {HTMLVideoElement} */ (document.getElementById("local_video")).srcObject = stream;
  for (const track of stream.getTracks()) {
    peer.connection.addTrack(track, stream);
  }
}

export async function receiveVideoOffer(message) {
  const { from: sender, sdp } = message;

  const peer = createPeerConnection(sender);

  const session = new RTCSessionDescription(sdp);
  await peer.connection.setRemoteDescription(session);

  let stream;
  try {
    stream = await navigator.mediaDevices.getUserMedia(mediaConstraints);
  } catch (error) {
    console.log(error);
    return;
  }

  /** @type {HTMLVideoElement} */ (document.getElementById("local_video")).srcObject = stream;

  for (const track of stream.getTracks()) {
    peer.connection.addTrack(track, stream);
  }

  try {
    const answer = await peer.connection.createAnswer();
    await peer.connection.setLocalDescription(answer);
  } catch (error) {
    console.log(error);
    return;
  }

  sendToServer({
    type: messages.VIDEO_ANSWER,
    to: sender,
    sdp: peer.connection.localDescription,
  });
}

export async function receiveICECandidate(message) {
  const candidate = new RTCIceCandidate(message.sdp);

  const peer = peers[message.from];
  if (!peer) return;

  try {
    await peer.connection.addIceCandidate(candidate);
  } catch (error) {
    console.log(error);
  }
}

/** @returns {Peer} */
export function createPeerConnection(peerName) {
  const peer = {
    connection: new RTCPeerConnection({
      iceServers: [
        {
          // Open source STUN service (https://www.stunprotocol.org/)
          urls: "stun:stun.stunprotocol.org",
        },
      ],
    }),
  };

  peer.connection.onicecandidate = (event) => handleICECandidate(peerName, event);
  peer.connection.ontrack = () => handleTrack(peer);
  peer.connection.onnegotiationneeded = () => handleNegotiationNeeded(peer, peerName);
  peer.connection.oniceconnectionstatechange = () => handleICEConnectionStateChange(peer);
  peer.connection.onicegatheringstatechange = () => handleICEGatheringStateChange(peer);
  peer.connection.onsignalingstatechange = () => handleSignalingStateChange(peer);

  peers[peerName] = peer;

  return peer;
}

/**
 * @param {string} peerName
 * @param {RTCPeerConnectionIceEvent} event
 */
function handleICECandidate(peerName, event) {
  if (event.candidate) {
    sendToServer({
      type: messages.ICE_CANDIDATE,
      to: peerName,
      sdp: event.candidate,
    });
  }
}

/** @param {Peer} peer */
async function handleNegotiationNeeded(peer, peerName) {
  try {
    const offer = await peer.connection.createOffer();
    await peer.connection.setLocalDescription(offer);
  } catch (error) {
    console.log(error);
    return;
  }

  sendToServer({
    type: messages.VIDEO_OFFER,
    to: peerName,
    sdp: peer.connection.localDescription,
  });
}

function handleTrack(peer) {}

function handleICEConnectionStateChange(peer) {}

function handleICEGatheringStateChange(peer) {}

function handleSignalingStateChange(peer) {}
