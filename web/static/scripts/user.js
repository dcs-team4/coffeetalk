/** @type {string} */
let username;

/** @returns {{ name: string, ok: true } | { ok: false }} */
export function getUsername() {
  if (username === undefined || username == "") {
    return { ok: false };
  }

  return { name: username, ok: true };
}

/** @param {string} name */
export function setUsername(name) {
  username = name;
}
