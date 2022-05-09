import { env } from "../env.js";
import { getUsername } from "../user.js";
import { leaveSession } from "../session.js";
import {
  sendVideoOffer,
  receiveVideoOffer,
  receiveVideoAnswer,
  receiveICECandidate,
  closePeerConnection,
} from "./peers.js";
import {
  displayError,
  setParticipantCount,
  incrementParticipantCount,
  decrementParticipantCount,
} from "../dom.js";

/**
 * WebSocket connection to the WebRTC signaling server. Undefined if uninitialized.
 * @type {WebSocket | undefined}
 */
let socket;

/**
 * Whether the socket is open and ready to send on.
 * @type {boolean}
 */
export let socketOpen = false;

/**
 * Tries to establish socket connection to the WebRTC signaling server.
 * If successful, registers event listeners on the socket.
 */
export function connectWebRTCSocket() {
  const protocol = env.ENV === "production" ? "wss" : "ws";
  const serverURL = `${protocol}://${env.WEBRTC_HOST}:${env.WEBRTC_PORT}`;

  try {
    socket = new WebSocket(serverURL);
  } catch (error) {
    console.log(`Socket connection to WebRTC signaling server failed:\n${error}`);
    return;
  }

  socket.addEventListener("message", handleMessage);
  socket.addEventListener("close", handleClose);
  socket.addEventListener("open", () => {
    socketOpen = true;
    console.log("Successfully connected to WebRTC signaling server.");
  });
  socket.addEventListener("error", (event) => {
    console.log(`Error in socket connection to WebRTC signaling server:\n${event}`);
  });
}

/**
 * Message types configured on the WebRTC server.
 * Should be updated if the server configuration changes.
 */
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

/**
 * Serializes the given message object to JSON, and sends it to the WebRTC signaling server.
 * @param {any} message
 */
export function sendWebRTCMessage(message) {
  if (!socket || !socketOpen) {
    console.log("Sending to server failed: socket uninitialized.");
    return;
  }

  console.log("Sending to server:", message);
  const serialized = JSON.stringify(message);
  socket.send(serialized);
}

/**
 * Handles the incoming WebSocket message.
 * @param {MessageEvent<any>} event
 */
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
      closePeerConnection(message.username);
      break;
    case messages.ERROR:
      console.log(`Error received from WebRTC server:\n${message.errorMessage}`);
      break;
    default:
      console.log(`Unrecognized message type: ${message.type}`);
  }
}

/** On socket connection close: leaves the stream, and reloads the app. */
function handleClose() {
  socketOpen = false;
  leaveSession();
  displayError("Connection lost. Reloading...");
  setTimeout(() => location.reload(), 5000);
}
