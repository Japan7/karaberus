import { platform, type Platform } from "@tauri-apps/plugin-os";
import { createEffect } from "solid-js";
import { removeSessionToken } from "../utils/session";
import {
  buildKaraberusUrl,
  getTauriUrl,
  IS_TAURI_DIST_BUILD,
} from "../utils/tauri";

export default function Logout() {
  createEffect(() => {
    removeSessionToken();
    if (IS_TAURI_DIST_BUILD) {
      const url = buildKaraberusUrl("/logout");
      url.searchParams.set("platform", platform());
      location.href = url.toString();
    } else {
      const params = new URLSearchParams(location.search);
      const platform = params.get("platform") as Platform | null;
      location.href = platform ? getTauriUrl(platform).toString() : "/";
    }
  });

  return null;
}
