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
import "./office.js";

// Initializes socket connection with WebRTC signaling server.
connectWebRTCSocket();

// Shows the users's own webcam.
streamLocalVideo();

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
