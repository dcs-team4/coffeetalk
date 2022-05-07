import { messages, sendToServer } from "./socket.js";
import { getUsername, setUsername } from "./user.js";
import { connectMQTT, startQuiz } from "./mqtt.js";
import { leaveCall } from "./webrtc.js";

/** DOM elements fetched on page load. */
export const DOM = {
  participantCount: () => /** @type {HTMLElement} */ (document.getElementById("participant-count")),
  errorField: () => /** @type {HTMLElement} */ (document.getElementById("error-field")),
  loginBar: () => /** @type {?HTMLElement} */ document.getElementById("login-bar"),
  usernameInput: () => /** @type {?HTMLInputElement} */ (document.getElementById("username-input")),
  joinCallButton: () => /** @type {?HTMLElement} */ document.getElementById("join-call"),
  quizBar: () => /** @type {HTMLElement} */ (document.getElementById("quiz-bar")),
  startQuizButton: () => /** @type {HTMLElement} */ (document.getElementById("start-quiz")),
  quizTitle: () => /** @type {HTMLElement} */ (document.getElementById("quiz-title")),
  quizQuestionTitle: () =>
    /** @type {HTMLElement} */ (document.getElementById("quiz-question-title")),
  quizQuestion: () => /** @type {HTMLElement} */ (document.getElementById("quiz-question")),
  quizAnswerTitle: () => /** @type {HTMLElement} */ (document.getElementById("quiz-answer-title")),
  quizAnswer: () => /** @type {HTMLElement} */ (document.getElementById("quiz-answer")),
  leaveStreamBar: () => /** @type {?HTMLElement} */ (document.getElementById("leave-stream-bar")),
  leaveStreamButton: () => /** @type {?HTMLElement} */ (document.getElementById("leave-stream")),
  videoContainer: () => /** @type {HTMLElement} */ (document.getElementById("video-container")),
  localVideo: () => /** @type {HTMLVideoElement} */ (document.getElementById("local-video")),
  localVideoName: () => /** @type {HTMLElement} */ (document.getElementById("local-video-name")),
};

registerListeners();

function registerListeners() {
  DOM.usernameInput()?.addEventListener("change", (event) => {
    setUsername(/** @type {HTMLInputElement} */ (event.target).value);
  });

  DOM.joinCallButton()?.addEventListener("click", joinCall);

  DOM.leaveStreamButton()?.addEventListener("click", leaveCall);

  DOM.startQuizButton().addEventListener("click", startQuiz);
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

  DOM.videoContainer()?.appendChild(container);

  return [video, container];
}

export function joinCall() {
  const user = getUsername();
  if (!user.ok) {
    DOM.errorField().innerText = "Invalid username";
    return;
  }
  DOM.errorField().innerText = "";

  sendToServer({
    type: messages.JOIN_STREAM,
    username: user.name,
  });

  connectMQTT();

  incrementParticipantCount();

  DOM.loginBar()?.classList.add("hide");
  DOM.quizBar().classList.remove("hide");
  DOM.leaveStreamBar()?.classList.remove("hide");
}

/** @param {number} value */
export function setParticipantCount(value) {
  DOM.participantCount().innerText = value.toString();
}

export function incrementParticipantCount() {
  const count = DOM.participantCount();
  const value = parseInt(count.innerText) || 0;
  count.innerText = (value + 1).toString();
}

export function decrementParticipantCount() {
  const count = DOM.participantCount();
  const value = parseInt(count.innerText);
  if (!value) {
    return;
  }

  count.innerText = (value - 1).toString();
}
