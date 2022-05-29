import { env } from "../env.js";
import { inStream, leaveStream } from "../stream.js";
import {
  handleNewPeer,
  receivePeerOffer,
  receivePeerAnswer,
  receiveICECandidate,
  closePeerConnection,
} from "./peers.js";
import { displayError, setPeerCount, incrementPeerCount, decrementPeerCount } from "../dom.js";

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
 * Message types configured on the WebRTC signaling server.
 * Should be updated if the server configuration changes.
 * @type {messages.MessageTypes}
 */
export const messages = {
  PEER_ANSWER: "peer-answer",
  PEER_OFFER: "peer-offer",
  ICE_CANDIDATE: "ice-candidate",
  JOIN_PEERS: "join-peers",
  LEAVE_PEERS: "leave-peers",
  PEER_JOINED: "peer-joined",
  PEER_LEFT: "peer-left",
  CONNECTION_SUCCESS: "connection-success",
  ERROR: "error",
};

/**
 * Serializes the given message to JSON, and sends it to the WebRTC signaling server.
 * @param {messages.SendableMessage} message
 */
export function sendWebRTCMessage(message) {
  if (!socket || !socketOpen) {
    console.log("Failed to send to WebRTC server: socket uninitialized.");
    return;
  }

  console.log("Sending to WebRTC server:", message);
  const serialized = JSON.stringify(message);
  socket.send(serialized);
}

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
    console.log("Socket connection to WebRTC server failed:", error);
    return;
  }

  socket.addEventListener("message", handleMessage);
  socket.addEventListener("close", handleClose);
  socket.addEventListener("open", () => {
    socketOpen = true;
    console.log("Successfully connected to WebRTC signaling server.");
  });
  socket.addEventListener("error", (event) => {
    console.log("Error in socket connection to WebRTC server:", event);
  });
}

/**
 * Handles the incoming WebSocket message.
 * @param {MessageEvent} event
 */
function handleMessage(event) {
  /** @type {messages.ReceivableMessage} */
  const message = JSON.parse(event.data);

  console.log("Received from WebRTC server:", message);

  switch (message.type) {
    case messages.PEER_OFFER:
      receivePeerOffer(message.from, message.data);
      break;
    case messages.PEER_ANSWER:
      receivePeerAnswer(message.from, message.data);
      break;
    case messages.ICE_CANDIDATE:
      receiveICECandidate(message.from, message.data);
      break;
    case messages.PEER_JOINED:
      incrementPeerCount();
      if (inStream) handleNewPeer(message.username);
      break;
    case messages.PEER_LEFT:
      decrementPeerCount();
      closePeerConnection(message.username);
      break;
    case messages.CONNECTION_SUCCESS:
      setPeerCount(message.peerCount);
      break;
    case messages.ERROR:
      console.log("Error from WebRTC server:", message.errorMessage);
      break;
    default: {
      console.log("WebRTC message not recognized.");
    }
  }
}

/** On socket connection close: leaves the stream, and reloads the app. */
function handleClose() {
  socketOpen = false;
  leaveStream();
  displayError("Connection lost. Reloading...");
  setTimeout(() => location.reload(), 5000);
}
