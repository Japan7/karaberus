export const isTauriBuild = import.meta.env.MODE.startsWith("tauri");
export const isTauriDevBuild = import.meta.env.MODE === "tauri-dev";
export const isTauriDistBuild = import.meta.env.MODE === "tauri-dist";

export const RELEASE_URL =
  "https://github.com/Japan7/karaberus/releases/latest";

export const REMOTE_DESKTOP_URL = `${import.meta.env.VITE_KARABERUS_URL}/desktop`;
