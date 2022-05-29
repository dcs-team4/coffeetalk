import { env } from "./env.js";
import "./user.js";
import "./stream.js";
import "./mqtt.js";
import { registerListeners } from "./dom.js";
import { connectWebRTCSocket } from "./webrtc/socket.js";
import { streamLocalVideo } from "./webrtc/peers.js";
import { initializeOffice, detectMotion } from "./office.js";

// Registers event listeners on DOM elements.
registerListeners();

// Initializes socket connection with WebRTC signaling server.
connectWebRTCSocket();

// Shows the users's own webcam.
streamLocalVideo();

// Runs office script if this is an office client.
if (env.CLIENT_TYPE === "office") {
  initializeOffice();
  detectMotion();
}
