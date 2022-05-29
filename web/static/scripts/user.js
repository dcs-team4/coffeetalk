import { DOM } from "./dom.js";

/**
 * Username used for stream ID and message passing. Undefined until user logs in.
 * @type {string | undefined}
 */
let username;

/**
 * Tries to get the user's username, if it's set.
 * @returns {{ name: string, ok: true } | { ok: false }}
 */
export function getUsername() {
  if (!username) {
    return { ok: false };
  }

  return { name: username, ok: true };
}

/**
 * Sets username to the given name.
 * @param {string} name
 */
export function setUsername(name) {
  username = name;
}

/**
 * Logs the user in with their selected username, and returns it.
 * If username is not selected, displays error message and returns `ok=false`.
 * @returns {{ name: string, ok: true } | { ok: false }}
 */
export function login() {
  const user = getUsername();

  if (user.ok) {
    DOM.errorField().innerText = "";
  } else {
    DOM.errorField().innerText = "Invalid username";
  }

  return user;
}
