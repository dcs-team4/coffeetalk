export function getUsername() {
  if (username) {
    return { username, ok: true };
  } else {
    return { username: undefined, ok: false };
  }
}
