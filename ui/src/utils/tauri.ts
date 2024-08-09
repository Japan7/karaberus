import type { Platform } from "@tauri-apps/plugin-os";

export const IS_TAURI_BUILD = import.meta.env.MODE.startsWith("tauri");
export const IS_TAURI_DEV_BUILD = import.meta.env.MODE === "tauri-dev";
export const IS_TAURI_DIST_BUILD = import.meta.env.MODE === "tauri-dist";

export const RELEASE_URL =
  "https://github.com/Japan7/karaberus/releases/latest";

export function getTauriUrl(platform: Platform) {
  return platform === "windows"
    ? "https://tauri.localhost"
    : "tauri://localhost";
}
