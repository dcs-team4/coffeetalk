//@ts-check

let username;

export function getUsername() {
  if (username) {
    return { username, ok: true };
  } else {
    return { username: undefined, ok: false };
  }
}

export function login(name) {
  username = name;
}

export function logout() {
  username = undefined;
}
