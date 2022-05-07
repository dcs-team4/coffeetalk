env = {
  ENV: "development",
  CLIENT_TYPE: "home",
  WEBRTC_HOST: "localhost",
  WEBRTC_PORT: "8000",
  MQTT_HOST: "localhost",
  MQTT_PORT: "1882",
  ...(env ?? {}),
};

import "./dom.js";
import "./user.js";
import { connectSocket } from "./socket.js";
import { streamLocalVideo } from "./webrtc.js";
import "./mqtt.js";
import "./office.js";

connectSocket();
streamLocalVideo();
