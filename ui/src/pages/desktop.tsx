import { createEffect } from "solid-js";
import { getSessionToken } from "../utils/session";

export default function Desktop() {
  createEffect(() => {
    const token = getSessionToken();
    location.href = "tauri://localhost?token=" + token;
  });

  return null;
}
