import { listen } from "@tauri-apps/api/event";
import type { Platform } from "@tauri-apps/plugin-os";

export const IS_TAURI_BUILD = import.meta.env.MODE.startsWith("tauri");
export const IS_TAURI_DEV_BUILD = import.meta.env.MODE === "tauri-dev";
export const IS_TAURI_DIST_BUILD = import.meta.env.MODE === "tauri-dist";

export const RELEASE_URL =
  "https://github.com/Japan7/karaberus/releases/latest";

export function getTauriUrl(platform: Platform, pathname = "/") {
  return new URL(
    pathname,
    platform === "windows" ? "http://tauri.localhost" : "tauri://localhost",
  );
}

export function buildKaraberusUrl(pathname: string) {
  return new URL(pathname, import.meta.env.VITE_KARABERUS_URL);
}

export function buildRedirectUrl(href: string) {
  const url = buildKaraberusUrl("/redirect");
  url.searchParams.set("href", href);
  return url;
}

export function registerGlobalListeners() {
  listen<string>("mpv-stdout", (e) => console.debug(e.payload));
  listen<string>("mpv-stderr", (e) => console.warn(e.payload));
}
