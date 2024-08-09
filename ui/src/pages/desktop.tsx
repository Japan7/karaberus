import type { Platform } from "@tauri-apps/plugin-os";
import { createEffect } from "solid-js";
import { getSessionToken } from "../utils/session";
import { getTauriUrl } from "../utils/tauri";

export default function Desktop() {
  createEffect(() => {
    const params = new URLSearchParams(location.search);
    const platform = params.get("platform") as Platform;
    const url = new URL(getTauriUrl(platform));
    url.searchParams.set("token", getSessionToken()!);
    location.href = url.toString();
  });

  return null;
}
