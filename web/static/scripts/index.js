//@ts-check

import "./dom.js";
import "./user.js";
import { connectSocket } from "./socket.js";
import { streamLocalVideo } from "./webrtc.js";
import "./mqtt.js";

connectSocket();
streamLocalVideo();
