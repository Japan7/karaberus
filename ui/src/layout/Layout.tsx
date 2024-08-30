import { useNavigate, type RouteSectionProps } from "@solidjs/router";
import { createEffect, Show } from "solid-js";
import { Provider } from "../components/Context";
import { getSessionInfos } from "../utils/session";
import AuthHero from "./AuthHero";
import Navbar from "./Navbar";
import Sidebar from "./Sidebar";

export default function Layout(props: RouteSectionProps) {
  const navigate = useNavigate();

  let checkboxRef!: HTMLInputElement;

  const infos = getSessionInfos();

  createEffect(() => {
    const preloginHref = localStorage.getItem("prelogin_href");
    if (preloginHref) {
      localStorage.removeItem("prelogin_href");
      const url = new URL(preloginHref);
      if (url.search) {
        location.href = preloginHref;
      } else if (location.pathname !== url.pathname) {
        navigate(url.pathname);
      }
    }
  });

  return (
    <Provider>
      <Show when={infos} fallback={<AuthHero />}>
        <div class="drawer lg:drawer-open">
          <input
            id="main-drawer"
            type="checkbox"
            ref={checkboxRef}
            class="drawer-toggle"
          />

          <div class="drawer-content">
            <div class="mt-2 mx-2">
              <Navbar infos={infos!} />
            </div>
            <main class="container flex flex-col py-6 gap-y-2">
              {props.children}
            </main>
          </div>

          <div class="drawer-side">
            <label
              for="main-drawer"
              aria-label="close sidebar"
              class="drawer-overlay"
            />
            <Sidebar
              closeDrawer={() => {
                checkboxRef.checked = false;
              }}
            />
          </div>
        </div>
      </Show>
    </Provider>
  );
}
