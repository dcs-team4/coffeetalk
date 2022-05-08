import { getUsername } from "./user.js";
import {
  closeConnections,
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
  displayError,
} from "./dom.js";

/** @type {WebSocket | undefined} */
let socket;

export function connectSocket() {
  const protocol = env.ENV === "production" ? "wss" : "ws";
  const host = env.MQTT_HOST;
  const port = env.WEBRTC_PORT;
  const serverURL = `${protocol}://${host}:${port}`;

  try {
    socket = new WebSocket(serverURL);
    socket.addEventListener("message", handleMessage);
    socket.addEventListener("close", handleClose);
    socket.addEventListener("open", () => {
      console.log("Successfully connected to WebRTC signaling server.");
    });
    socket.addEventListener("error", (event) => {
      console.log(`Error in socket connection to WebRTC signaling server:\n${event}`);
    });
  } catch (error) {
    console.log(`Socket connection to WebRTC signaling server failed:\n${error}`);
  }
}

/** @param {any} message */
export function sendToServer(message) {
  if (!socket) {
    console.log("Sending to server failed: socket uninitialized.");
    return;
  }

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
      receiveICECandidate(message.from, message.sdp);
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

function handleClose() {
  closeConnections();
  setParticipantCount(0);
  displayError("Connection lost. Reloading...");
  setTimeout(() => location.reload(), 5000);
}
