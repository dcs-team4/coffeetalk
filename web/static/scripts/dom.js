import { setUsername } from "./user.js";
import { startQuiz } from "./mqtt.js";
import { joinSession, leaveSession } from "./session.js";

/**
 * Utility object wrapping functions for fetching DOM elements.
 * Elements typed as non-optional should be assumed to be present; those typed as optional depend
 * on how the web app is accessed (e.g. home vs. office client).
 */
export const DOM = Object.freeze({
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
  videos: () => /** @type {HTMLElement} */ (document.getElementById("videos")),
  localVideo: () => /** @type {HTMLVideoElement} */ (document.getElementById("local-video")),
  localVideoName: () => /** @type {HTMLElement} */ (document.getElementById("local-video-name")),
});

registerListeners();

/** Initializes event listeners on the appropriate DOM objects. */
function registerListeners() {
  DOM.usernameInput()?.addEventListener("change", (event) => {
    const name = /** @type {HTMLInputElement} */ (event.target).value;
    setUsername(name);
  });

  DOM.joinCallButton()?.addEventListener("click", joinSession);
  DOM.leaveStreamButton()?.addEventListener("click", leaveSession);
  DOM.startQuizButton().addEventListener("click", startQuiz);
}

/**
 * Creates a video element for a peer, wraps it in a container, and appends it to the stream view.
 * Adds the given `peerName` in a block under the video.
 * Returns the video element and the container surrounding it.
 * @param {string} peerName
 * @returns {{ video: HTMLVideoElement, container: HTMLElement }}
 */
export function createPeerVideoElement(peerName) {
  const container = document.createElement("div");

  const video = document.createElement("video");
  video.autoplay = true;
  container.appendChild(video);

  const name = document.createElement("div");
  name.innerText = peerName;
  name.classList.add("video-name");
  container.appendChild(name);

  DOM.videos()?.appendChild(container);

  return { video, container };
}

/**
 * Sets the participant count display to the given value.
 * @param {number} value
 */
export function setParticipantCount(value) {
  DOM.participantCount().innerText = value.toString();
}

/**
 * Increments the participant count display by one.
 */
export function incrementParticipantCount() {
  const count = DOM.participantCount();
  const value = parseInt(count.innerText) || 0;
  count.innerText = (value + 1).toString();
}

/**
 * Decrements the participant count display by one, unless it is missing/already 0.
 */
export function decrementParticipantCount() {
  const count = DOM.participantCount();
  const value = parseInt(count.innerText);
  if (!value) {
    return;
  }

  count.innerText = (value - 1).toString();
}

/** Shows DOM elements for the active stream view, and hides others. */
export function displayStream() {
  DOM.loginBar()?.classList.add("hide");
  DOM.quizBar().classList.remove("hide");
  DOM.leaveStreamBar()?.classList.remove("hide");
}

/** Shows DOM elements for the login view, and hides others. */
export function displayLogin() {
  DOM.leaveStreamBar()?.classList.add("hide");
  DOM.quizBar().classList.add("hide");
  DOM.loginBar()?.classList.remove("hide");
}

/**
 * Shows the given error message in the error display.
 * @param {string} errorMessage
 */
export function displayError(errorMessage) {
  DOM.errorField().innerText = errorMessage;
}

/** Clears the text of the error display. */
export function clearError() {
  DOM.errorField().innerText = "";
}
