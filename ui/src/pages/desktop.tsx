import { useParams } from "@solidjs/router";
import type { Platform } from "@tauri-apps/plugin-os";
import { createEffect } from "solid-js";
import { getSessionToken } from "../utils/session";

export default function Desktop() {
  const params = useParams();

  createEffect(() => {
    const token = getSessionToken();
    const platform = params.platform as Platform;
    const tauriUrl =
      platform === "windows" ? "https://tauri.localhost" : "tauri://localhost";
    location.href = tauriUrl + "?token=" + token;
  });

  return null;
}
