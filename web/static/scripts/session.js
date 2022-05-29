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
 * - Signals to the WebRTC server that we are ready to set up peer connections.
 * - Connects to the MQTT broker for quiz sessions.
 * - Displays the stream.
 */
export function joinSession() {
  const user = login();
  if (!user.ok) return;

  sendWebRTCMessage({ type: messages.JOIN_PEERS, username: user.name });
  connectMQTT();
  displayStream();
  incrementPeerCount();

  inSession = true;
  console.log("Joined session.");
}

/**
 * - Closes peer-to-peer connections, and alerts the WebRTC server that we are leaving.
 * - Disconnects from the MQTT broker.
 * - Returns to the login view.
 */
export function leaveSession() {
  closePeerConnections();
  sendWebRTCMessage({ type: messages.LEAVE_PEERS });
  disconnectMQTT();
  displayLogin();
  decrementPeerCount();

  inSession = false;
  console.log("Left session.");
}
