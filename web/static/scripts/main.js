import { env } from "./env.js";
import "./dom.js";
import "./user.js";
import "./session.js";
import "./mqtt.js";
import { connectWebRTCSocket } from "./webrtc/socket.js";
import { streamLocalVideo } from "./webrtc/peers.js";
import { initializeOffice, detectMotion } from "./office.js";

// Initializes socket connection with WebRTC signaling server.
connectWebRTCSocket();

// Shows the users's own webcam.
streamLocalVideo();

// Runs office script if this is an office client.
if (env.CLIENT_TYPE === "office") {
  initializeOffice();
  detectMotion();
}
