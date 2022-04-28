//@ts-check

import { connectSocket } from "./socket.js";
import "./webrtc.js";

connectSocket();

window.onload = () => {
  const hangupButton = document.querySelector("#hangup-button");
  hangupButton.addEventListener("click", hangupCall);
};

function hangupCall() {}
