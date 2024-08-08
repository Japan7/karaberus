export const IS_TAURI_BUILD = import.meta.env.MODE.startsWith("tauri");
export const IS_TAURI_DEV_BUILD = import.meta.env.MODE === "tauri-dev";
export const IS_TAURI_DIST_BUILD = import.meta.env.MODE === "tauri-dist";

export const RELEASE_URL =
  "https://github.com/Japan7/karaberus/releases/latest";

export const REMOTE_DESKTOP_URL = `${import.meta.env.VITE_KARABERUS_URL}/desktop`;
