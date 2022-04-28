import { videoAnswerMessage, videoOfferMessage } from "./messages.js";
import { sendToServer } from "./socket.js";

/** @type {{ [key: string]: RTCPeerConnection }} */
const peers = {};

const mediaConstraints = {
  audio: true,
  video: true,
};

export function sendVideoOffer(peerName) {
  const peer = createPeerConnection(peerName);

  let stream;
  try {
    stream = await navigator.mediaDevices.getUserMedia(mediaConstraints);
  } catch (error) {
    console.log(error);
    return;
  }

  document.getElementById("local_video").srcObject = stream;
  for (const track of stream.getTracks()) {
    peer.addTrack(track, stream);
  }
}

export function receiveVideoOffer(message) {
  const { from, sdp } = message;

  const peer = createPeerConnection(from);

  const session = new RTCSessionDescription(sdp);
  await peer.setRemoteDescription(session);

  let stream;
  try {
    const stream = await navigator.mediaDevices.getUserMedia(mediaConstraints);
  } catch (error) {
    console.log(error);
    return;
  }
  
  document.getElementById("local_video").srcObject = stream;
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

  const message = videoAnswerMessage(from, peer.localDescription);
  if (message) {
    sendToServer(message);
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

  peer.onicecandidate = () => handleICECandidate(peer);
  peer.ontrack = () => handleTrack(peer);
  peer.onnegotiationneeded = () => handleNegotiationNeeded(peer, peerName);
  peer.oniceconnectionstatechange = () => handleICEConnectionStateChange(peer);
  peer.onicegatheringstatechange = () => handleICEGatheringStateChange(peer);
  peer.onsignalingstatechange = () => handleSignalingStateChange(peer);

  peers[peerName] = peer;

  return peer;
}

function handleICECandidate(peer) {}

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

  const message = videoOfferMessage(peerName, peer.localDescription);
  if (message) {
    sendToServer(message);
  }
}

function handleICEConnectionStateChange(peer) {}

function handleICEGatheringStateChange(peer) {}

function handleSignalingStateChange(peer) {}
