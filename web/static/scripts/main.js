// Sets default environment variables, overwritten by env passed from server.
env = {
  ENV: "development",
  CLIENT_TYPE: "home",
  WEBRTC_HOST: "localhost",
  WEBRTC_PORT: "8000",
  MQTT_HOST: "localhost",
  MQTT_PORT: "1882",
  ...(env ?? {}),
};

import {
  displayStream,
  displayLogin,
  incrementParticipantCount,
  decrementParticipantCount,
} from "./dom.js";
import { login } from "./user.js";
import { connectWebRTCSocket, sendWebRTCMessage, messages } from "./webrtc/socket.js";
import { streamLocalVideo, closePeerConnections } from "./webrtc/peers.js";
import { connectMQTT, disconnectMQTT } from "./mqtt.js";

// Initializes socket connection with WebRTC signaling server.
connectWebRTCSocket();

// Shows the users's own webcam.
streamLocalVideo();

// Runs office script if this is an office client.
if (env.CLIENT_TYPE === "office") {
  runOfficeScript();
}

/** Loads the script for the office client, and runs it. */
async function runOfficeScript() {
  try {
    const { initializeOffice, detectMotion } = await import("./office.js");
    initializeOffice();
    detectMotion();
  } catch (error) {
    console.log("Failed to load office script:", error);
  }
}

/**
 * - Logs in with the user's provided username
 * - Alerts the WebRTC server that we are ready to join the stream
 * - Connects to the MQTT broker for quiz sessions
 * - Displays the stream
 */
export function joinCall() {
  const user = login();
  if (!user.ok) return;

  sendWebRTCMessage({ type: messages.JOIN_STREAM, username: user.name });
  connectMQTT();
  incrementParticipantCount();
  displayStream();
}

/**
 * - Closes peer-to-peer connections.
 * - Displays the stream
 * - Alerts the WebRTC signaling server that we want to leave the stream.
 * - Disconnects from the MQTT broker.
 */
export function leaveCall() {
  closePeerConnections();
  displayLogin();
  decrementParticipantCount();
  sendWebRTCMessage({ type: messages.LEAVE_STREAM });
  disconnectMQTT();
}
