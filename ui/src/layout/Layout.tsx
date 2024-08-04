import type { RouteSectionProps } from "@solidjs/router";
import { Show } from "solid-js";
import { getSessionInfos } from "../utils/session";
import AuthHero from "./AuthHero";
import Navbar from "./Navbar";
import Sidebar from "./Sidebar";

export default function Layout({ children }: RouteSectionProps) {
  const infos = getSessionInfos();

  return (
    <Show when={infos} fallback={<AuthHero />}>
      <div class="drawer lg:drawer-open">
        <input id="main-drawer" type="checkbox" class="drawer-toggle" />

        <div class="drawer-content">
          <div class="mt-2 mx-2">
            <Navbar infos={infos!} />
          </div>
          <main class="container flex flex-col py-6 gap-y-2">{children}</main>
        </div>

        <div class="drawer-side">
          <label
            for="main-drawer"
            aria-label="close sidebar"
            class="drawer-overlay"
          />
          <Sidebar />
        </div>
      </div>
    </Show>
  );
}
