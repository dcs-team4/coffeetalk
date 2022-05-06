//@ts-check

import {
  createPeerConnection,
  receiveICECandidate,
  receiveVideoOffer,
  removePeerConnection,
} from "./webrtc.js";

/** @type {WebSocket | undefined} */
let socket;

export function connectSocket() {
  const protocol = env?.ENV === "production" ? "wss" : "ws";
  const host = env?.MQTT_HOST ?? "localhost";
  const port = env?.WEBRTC_PORT ?? "8000";
  const serverURL = `${protocol}://${host}:${port}`;

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
  USER_JOINED: "user-joined",
  LEAVE_STREAM: "leave-stream",
  USER_LEFT: "user-left",
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
    case messages.USER_JOINED:
      createPeerConnection(message.username);
      break;
    case messages.USER_LEFT:
      removePeerConnection(message.username);
      break;
    case messages.ERROR:
      console.log(`Error received from WebRTC server: ${message.errorMessage}`);
      break;
    default:
      console.log(`Unrecognized message type: ${message.type}`);
  }
}
