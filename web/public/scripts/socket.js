//@ts-check

import { receiveICECandidate, receiveVideoOffer } from "./webrtc";

const socketPort = "8000";

/** @type {WebSocket | undefined} */
let socket;

export function connectSocket() {
  let protocol = "ws";
  if (location.protocol === "https:") {
    protocol = "wss";
  }
  const serverURL = `${protocol}://${location.hostname}:${socketPort}`;

  socket = new WebSocket(serverURL);
  socket.addEventListener("message", handleMessage);
}

export function sendToServer(message) {
  const serialized = JSON.stringify(message);
  socket.send(serialized);
}

export const messages = Object.freeze({
  VIDEO_OFFER: "video-offer",
  VIDEO_ANSWER: "video-answer",
  ICE_CANDIDATE: "new-ice-candidate",
  JOIN_STREAM: "join-stream",
  LEAVE_STREAM: "leave-stream",
  ERROR: "error",
});

/** @param {MessageEvent<any>} event */
function handleMessage(event) {
  const message = JSON.parse(event.data);

  switch (message.type) {
    case messages.VIDEO_OFFER:
    case messages.VIDEO_ANSWER:
      receiveVideoOffer(message);
      break;
    case messages.ICE_CANDIDATE:
      receiveICECandidate(message);
      break;
    case messages.ERROR:
      console.log(`Error received from WebRTC server: ${message.errorMessage}`);
    default:
      console.log(`Unrecognized message type: ${message.type}`);
  }
}
