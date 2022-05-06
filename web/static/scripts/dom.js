//@ts-check

import { messages, sendToServer } from "./socket.js";
import { getUsername, setUsername } from "./user.js";
import { connectMQTT, startQuiz } from "./mqtt.js";
import { leaveCall } from "./webrtc.js";

/** DOM elements fetched on page load. */
export const DOM = {
  participantCount: document.getElementById("participant-count"),
  loginBar: document.getElementById("login-bar"),
  loginError: document.getElementById("login-error"),
  usernameInput: /** @type {HTMLInputElement} */ (document.getElementById("username-input")),
  joinCallButton: document.getElementById("join-call"),
  quizBar: document.getElementById("quiz-bar"),
  startQuizButton: document.getElementById("start-quiz"),
  quizTitle: document.getElementById("quiz-title"),
  quizQuestionTitle: document.getElementById("quiz-question-title"),
  quizQuestion: document.getElementById("quiz-question"),
  quizAnswerTitle: document.getElementById("quiz-answer-title"),
  quizAnswer: document.getElementById("quiz-answer"),
  leaveStreamBar: document.getElementById("leave-stream-bar"),
  leaveStreamButton: document.getElementById("leave-stream"),
  videoContainer: document.getElementById("video-container"),
  localVideo: /** @type {HTMLVideoElement} */ (document.getElementById("local-video")),
};

registerListeners();

function registerListeners() {
  DOM.usernameInput.addEventListener("change", (event) => {
    setUsername(/** @type {HTMLInputElement} */ (event.target).value);
  });

  DOM.joinCallButton.addEventListener("click", handleJoinCall);

  DOM.startQuizButton.addEventListener("click", startQuiz);

  DOM.leaveStreamButton.addEventListener("click", leaveCall);
}

/** @param {string} peerName, @returns {[HTMLVideoElement, HTMLElement]} */
export function createPeerVideoElement(peerName) {
  const container = document.createElement("div");

  const video = document.createElement("video");
  video.autoplay = true;
  container.appendChild(video);

  const nameEl = document.createElement("div");
  nameEl.innerText = peerName;
  nameEl.classList.add("video-name");
  container.appendChild(nameEl);

  DOM.videoContainer.appendChild(container);

  return [video, container];
}

function handleJoinCall() {
  const user = getUsername();
  if (!user.ok) {
    DOM.loginError.innerText = "Invalid username";
    return;
  }
  DOM.loginError.innerText = "";

  sendToServer({
    type: messages.JOIN_STREAM,
    username: user.name,
  });

  connectMQTT();

  incrementParticipantCount();

  DOM.loginBar.classList.add("hide");
  DOM.quizBar.classList.remove("hide");
  DOM.leaveStreamBar.classList.remove("hide");
}

/** @param {number} value */
export function setParticipantCount(value) {
  DOM.participantCount.innerText = value.toString();
}

export function incrementParticipantCount() {
  const value = parseInt(DOM.participantCount.innerText) || 0;
  DOM.participantCount.innerText = (value + 1).toString();
}

export function decrementParticipantCount() {
  const value = parseInt(DOM.participantCount.innerText);
  if (!value) {
    return;
  }

  DOM.participantCount.innerText = (value - 1).toString();
}
