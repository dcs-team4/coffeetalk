const connectAddress = "ws://localhost:8000/connect";

const socket = new WebSocket(connectAddress);

export function sendToServer(message) {
  const serialized = JSON.stringify(message);
  socket.send(serialized);
}
