import { isTauri } from "@tauri-apps/api/core";
import { listen, type UnlistenFn } from "@tauri-apps/api/event";
import { onOpenUrl } from "@tauri-apps/plugin-deep-link";
import { open } from "@tauri-apps/plugin-shell";
import { createEffect, onCleanup } from "solid-js";
import { apiUrl } from "../utils/karaberus-client";
import { setSessionToken } from "../utils/session";

export default function AuthHero() {
  let openUrlListener: UnlistenFn | undefined;
  let deepLinkListener: UnlistenFn | undefined;

  createEffect(async () => {
    if (isTauri()) {
      // macOS/mobile
      openUrlListener = await onOpenUrl((urls) => {
        urls.forEach(readDeepLink);
      });

      // Linux/Windows
      deepLinkListener = await listen<string>("deep-link", (e) => {
        readDeepLink(e.payload);
      });
    }
  });

  const readDeepLink = (url: string) => {
    const urlObj = new URL(url);
    const token = urlObj.searchParams.get("token");
    if (token) {
      setSessionToken(token);
      location.reload();
    }
  };

  const openBrowserConnect = () => {
    open(`${import.meta.env.VITE_KARABERUS_URL}/desktop`);
  };

  const redirectToLogin = () => {
    localStorage.setItem("prelogin_path", location.pathname);
    location.href = apiUrl("api/oidc/login");
  };

  onCleanup(() => {
    openUrlListener?.();
    deepLinkListener?.();
  });

  return (
    <div class="hero bg-base-200 min-h-screen [background-image:url(https://cdn.myanimelist.net/s/common/uploaded_files/1445139435-b6abfa181eae79d82e5eb41cf52bb72f.jpeg)]">
      <div class="hero-overlay bg-opacity-60"></div>
      <div class="hero-content text-neutral-content flex-col lg:flex-row-reverse lg:gap-x-12">
        <div class="text-center lg:text-left">
          <h1 class="text-5xl font-bold">Karaberus</h1>
          <p class="py-6">wow such empty</p>
        </div>
        <div class="card bg-base-100 bg-opacity-60 w-full max-w-sm shrink-0 shadow-2xl">
          <div class="card-body">
            <div class="form-control">
              <button
                onclick={isTauri() ? openBrowserConnect : redirectToLogin}
                class="btn btn-primary"
              >
                {isTauri()
                  ? "Open browser for login"
                  : "Login with OpenID Connect"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
