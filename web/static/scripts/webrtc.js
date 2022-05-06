//@ts-check

import { messages, sendToServer } from "./socket.js";
import { createPeerVideoElement, decrementParticipantCount, DOM } from "./dom.js";
import { getUsername } from "./user.js";

/**
 * @typedef {Object} Peer
 * @property {RTCPeerConnection} connection
 * @property {HTMLElement} [videoContainer]
 * @property {HTMLVideoElement} [video]
 */

/** @type {{ [username: string]: Peer }} */
const peers = {};

/** @type {MediaStream} */
let localStream;

const mediaConstraints = {
  audio: true,
  video: true,
};

export async function streamLocalVideo() {
  try {
    localStream = await navigator.mediaDevices.getUserMedia(mediaConstraints);
  } catch (error) {
    console.log(error);
    return;
  }

  DOM.localVideo.srcObject = localStream;
}

export function leaveCall() {
  for (const peer of Object.values(peers)) {
    peer.connection.close();
    DOM.videoContainer.removeChild(peer.videoContainer);
  }

  for (const peerName in peers) {
    delete peers[peerName];
  }

  DOM.leaveStreamBar.classList.add("hide");
  DOM.quizBar.classList.add("hide");
  DOM.loginBar.classList.remove("hide");

  decrementParticipantCount();

  sendToServer({ type: messages.LEAVE_STREAM });
}

/** @param {string} peerName */
export async function sendVideoOffer(peerName) {
  const peer = createPeerConnection(peerName);

  if (!localStream) return;

  for (const track of localStream.getTracks()) {
    peer.connection.addTrack(track, localStream);
  }
}

/** @param {string} sender, @param {any} sdp */
export async function receiveVideoAnswer(sender, sdp) {
  if (!peers.hasOwnProperty(sender)) {
    console.log(`Unrecognized peer: ${sender}`);
    return;
  }

  const peer = peers[sender];

  const session = new RTCSessionDescription(sdp);
  await peer.connection.setRemoteDescription(session);
}

/** @param {string} sender, @param {any} sdp */
export async function receiveVideoOffer(sender, sdp) {
  const peer = createPeerConnection(sender);

  const session = new RTCSessionDescription(sdp);
  await peer.connection.setRemoteDescription(session);

  if (localStream) {
    for (const track of localStream.getTracks()) {
      peer.connection.addTrack(track, localStream);
    }
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

/** @param {string} peerName */
export function removePeerConnection(peerName) {
  if (!peers.hasOwnProperty(peerName)) return;
  const peer = peers[peerName];

  peer.connection.close();
  if (peer.videoContainer) DOM.videoContainer.removeChild(peer.videoContainer);
  delete peers[peerName];
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

  peer.connection.onicecandidate = (event) => handleICECandidate(event, peerName);
  peer.connection.ontrack = (event) => handleTrack(event, peer, peerName);
  peer.connection.onnegotiationneeded = () => handleNegotiationNeeded(peer, peerName);
  peer.connection.oniceconnectionstatechange = () => handleICEConnectionStateChange(peer);
  peer.connection.onicegatheringstatechange = () => handleICEGatheringStateChange(peer);
  peer.connection.onsignalingstatechange = () => handleSignalingStateChange(peer);

  peers[peerName] = peer;

  return peer;
}

/** @param {RTCPeerConnectionIceEvent} event, @param {string} peerName */
function handleICECandidate(event, peerName) {
  if (event.candidate) {
    sendToServer({
      type: messages.ICE_CANDIDATE,
      to: peerName,
      sdp: event.candidate,
    });
  }
}

/** @param {Peer} peer, @param {string} peerName */
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

/** @param {RTCTrackEvent} event, @param {Peer} peer, @param {string} peerName */
function handleTrack(event, peer, peerName) {
  if (!peer.video) {
    const [video, container] = createPeerVideoElement(peerName);
    peer.video = video;
    peer.videoContainer = container;
  }

  peer.video.srcObject = event.streams[0];
}

function handleICEConnectionStateChange(peer) {}

function handleICEGatheringStateChange(peer) {}

function handleSignalingStateChange(peer) {}
