import { displayLogin, displayStream, incrementPeerCount, decrementPeerCount } from "./dom.js";
import { connectMQTT, disconnectMQTT } from "./mqtt.js";
import { login } from "./user.js";
import { closePeerConnections } from "./webrtc/peers.js";
import { sendWebRTCMessage, messages } from "./webrtc/socket.js";

/** Whether the user has joined the CoffeeTalk session. */
export let inSession = false;

/**
 * If not already in the session:
 * - Logs in with the user's provided username.
 * - Alerts the WebRTC server that we are ready to join the stream.
 * - Connects to the MQTT broker for quiz sessions.
 * - Displays the stream.
 */
export function joinSession() {
  const user = login();
  if (!user.ok) return;

  sendWebRTCMessage({ type: messages.JOIN_STREAM, username: user.name });
  connectMQTT();
  displayStream();
  incrementPeerCount();

  inSession = true;
  console.log("Joined session.");
}

/**
 * - Closes peer-to-peer connections.
 * - Returns to the login view.
 * - Alerts the WebRTC signaling server that we want to leave the stream.
 * - Disconnects from the MQTT broker.
 */
export function leaveSession() {
  closePeerConnections();
  sendWebRTCMessage({ type: messages.LEAVE_STREAM });
  disconnectMQTT();
  displayLogin();
  decrementPeerCount();

  inSession = false;
  console.log("Left session.");
}
