import { getUsername } from "./user";

const messageTypes = Object.freeze({
  videoOffer: "video-offer",
  videoAnswer: "video-answer",
  iceCandidate: "new-ice-candidate",
  joinStream: "join-stream",
  error: "error",
});

export function videoOfferMessage(to, sdp) {
  const { username, ok } = getUsername();
  if (!ok) return undefined;

  return {
    type: messageTypes.videoOffer,
    from: username,
    to,
    sdp,
  };
}

export function videoAnswerMessage(to, sdp) {
  const { username, ok } = getUsername();
  if (!ok) return undefined;

  return {
    type: messageTypes.videoOffer,
    from: username,
    to,
    sdp,
  };
}
