//@ts-check

import * as messages from "./messages.js";
import { sendToServer } from "./socket.js";

/** @type {{ [key: string]: RTCPeerConnection }} */
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
    peer.addTrack(track, stream);
  }
}

export async function receiveVideoOffer(message) {
  const { from: sender, sdp } = message;

  const peer = createPeerConnection(sender);

  const session = new RTCSessionDescription(sdp);
  await peer.setRemoteDescription(session);

  let stream;
  try {
    stream = await navigator.mediaDevices.getUserMedia(mediaConstraints);
  } catch (error) {
    console.log(error);
    return;
  }

  /** @type {HTMLVideoElement} */ (document.getElementById("local_video")).srcObject = stream;

  for (const track of stream.getTracks()) {
    peer.addTrack(track, stream);
  }

  try {
    const answer = await peer.createAnswer();
    await peer.setLocalDescription(answer);
  } catch (error) {
    console.log(error);
    return;
  }

  sendToServer({
    type: messages.VIDEO_ANSWER,
    to: sender,
    sdp: peer.localDescription,
  });
}

export async function receiveICECandidate(message) {
  const candidate = new RTCIceCandidate(message.sdp);

  const peer = peers[message.from];
  if (!peer) return;

  try {
    await peer.addIceCandidate(candidate);
  } catch (error) {
    console.log(error);
  }
}

export function createPeerConnection(peerName) {
  const peer = new RTCPeerConnection({
    iceServers: [
      {
        // Open source STUN service (https://www.stunprotocol.org/)
        urls: "stun:stun.stunprotocol.org",
      },
    ],
  });

  peer.onicecandidate = (event) => handleICECandidate(peer, peerName, event);
  peer.ontrack = () => handleTrack(peer);
  peer.onnegotiationneeded = () => handleNegotiationNeeded(peer, peerName);
  peer.oniceconnectionstatechange = () => handleICEConnectionStateChange(peer);
  peer.onicegatheringstatechange = () => handleICEGatheringStateChange(peer);
  peer.onsignalingstatechange = () => handleSignalingStateChange(peer);

  peers[peerName] = peer;

  return peer;
}

/**
 * @param {RTCPeerConnection} peer
 * @param {string} peerName
 * @param {RTCPeerConnectionIceEvent} event
 */
function handleICECandidate(peer, peerName, event) {
  if (event.candidate) {
    sendToServer({
      type: messages.ICE_CANDIDATE,
      to: peerName,
      sdp: event.candidate,
    });
  }
}

function handleTrack(peer) {}

/** @param {RTCPeerConnection} peer */
async function handleNegotiationNeeded(peer, peerName) {
  try {
    const offer = await peer.createOffer();
    await peer.setLocalDescription(offer);
  } catch (error) {
    console.log(error);
    return;
  }

  sendToServer({
    type: messages.VIDEO_OFFER,
    to: peerName,
    sdp: peer.localDescription,
  });
}

function handleICEConnectionStateChange(peer) {}

function handleICEGatheringStateChange(peer) {}

function handleSignalingStateChange(peer) {}
