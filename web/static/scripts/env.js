// Sets default environment variables, overwritten by env passed from server.
export const env = {
  ENV: "development",
  CLIENT_TYPE: "home",
  WEBRTC_HOST: "localhost",
  WEBRTC_PORT: "8000",
  MQTT_HOST: "localhost",
  MQTT_PORT: "1882",
  ...(serverEnv ?? {}),
};
