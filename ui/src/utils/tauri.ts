export const isTauri = () =>
  "__TAURI__" in window || "__TAURI_INTERNALS__" in window; // v1 || v2

export const RELEASE_URL =
  "https://github.com/Japan7/karaberus/releases/latest";
