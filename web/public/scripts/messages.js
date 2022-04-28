import { getUsername } from "./user";

const messageTypes = Object.freeze({
  videoOffer: "video-offer",
  videoAnswer: "video-answer",
  iceCandidate: "new-ice-candidate",
  joinStream: "join-stream",
  error: "error",
});

function videoExchangeMessage(message) {
  const { username, ok } = getUsername();
  if (!ok) return undefined;

  return {
    ...message,
    from: username,
  };
}

export function videoOfferMessage(to, sdp) {
  return videoExchangeMessage({ type: messageTypes.videoOffer, to, sdp });
}

export function videoAnswerMessage(to, sdp) {
  return videoExchangeMessage({ type: messageTypes.videoAnswer, to, sdp });
}

export function iceCandidateMessage(to, sdp) {
  return videoExchangeMessage({ type: messageTypes.iceCandidate, to, sdp });
}
