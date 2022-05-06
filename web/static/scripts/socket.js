//@ts-check

import { getUsername } from "./user.js";
import {
  receiveICECandidate,
  receiveVideoAnswer,
  receiveVideoOffer,
  removePeerConnection,
  sendVideoOffer,
} from "./webrtc.js";
import {
  setParticipantCount,
  incrementParticipantCount,
  decrementParticipantCount,
} from "./dom.js";

/** @type {WebSocket | undefined} */
let socket;

export function connectSocket() {
  const protocol = env?.ENV === "production" ? "wss" : "ws";
  const host = env?.MQTT_HOST ?? "localhost";
  const port = env?.WEBRTC_PORT ?? "8000";
  const serverURL = `${protocol}://${host}:${port}`;

  socket = new WebSocket(serverURL);
  socket.addEventListener("message", handleMessage);
  console.log("Successfully connected to WebRTC signaling server.");
}

export function sendToServer(message) {
  console.log("Sending to server:", message);
  const serialized = JSON.stringify(message);
  socket.send(serialized);
}

export const messages = Object.freeze({
  VIDEO_OFFER: "video-offer",
  VIDEO_ANSWER: "video-answer",
  ICE_CANDIDATE: "new-ice-candidate",
  JOIN_STREAM: "join-stream",
  CONNECTION_SUCCESS: "connection-success",
  USER_JOINED: "user-joined",
  LEAVE_STREAM: "leave-stream",
  USER_LEFT: "user-left",
  ERROR: "error",
});

/** @param {MessageEvent<any>} event */
function handleMessage(event) {
  const message = JSON.parse(event.data);
  console.log("Socket message received:", message);

  switch (message.type) {
    case messages.VIDEO_OFFER:
      receiveVideoOffer(message.from, message.sdp);
      break;
    case messages.VIDEO_ANSWER:
      receiveVideoAnswer(message.from, message.sdp);
      break;
    case messages.ICE_CANDIDATE:
      receiveICECandidate(message);
      break;
    case messages.CONNECTION_SUCCESS:
      setParticipantCount(message.participantCount);
      break;
    case messages.USER_JOINED:
      incrementParticipantCount();

      const user = getUsername();
      if (user.ok) sendVideoOffer(message.username);
      break;
    case messages.USER_LEFT:
      decrementParticipantCount();
      removePeerConnection(message.username);
      break;
    case messages.ERROR:
      console.log(`Error received from WebRTC server: ${message.errorMessage}`);
      break;
    default:
      console.log(`Unrecognized message type: ${message.type}`);
  }
}
