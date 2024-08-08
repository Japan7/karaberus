import { platform } from "@tauri-apps/plugin-os";
import { createEffect } from "solid-js";
import { apiUrl } from "../utils/karaberus-client";
import { setSessionToken } from "../utils/session";
import { IS_TAURI_DIST_BUILD, REMOTE_DESKTOP_URL } from "../utils/tauri";

export default function AuthHero() {
  createEffect(() => {
    if (IS_TAURI_DIST_BUILD) {
      const token = new URLSearchParams(location.search).get("token");
      if (token) {
        setSessionToken(token);
        location.href = "/";
      } else {
        const platformName = platform();
        location.href = REMOTE_DESKTOP_URL + "?platform=" + platformName;
      }
    }
  });

  const redirectToLogin = () => {
    localStorage.setItem("prelogin_path", location.pathname);
    location.href = apiUrl("api/oidc/login");
  };

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
              <button onclick={redirectToLogin} class="btn btn-primary">
                Login with OpenID Connect
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
