import type { RouteSectionProps } from "@solidjs/router";
import { getSessionInfos } from "../utils/session";
import AuthHero from "./AuthHero";
import Navbar from "./Navbar";
import Sidebar from "./Sidebar";

export default function Layout({ children }: RouteSectionProps) {
  const infos = getSessionInfos();

  return !infos ? (
    <AuthHero />
  ) : (
    <div class="drawer lg:drawer-open">
      <input id="main-drawer" type="checkbox" class="drawer-toggle" />

      <div class="drawer-content">
        <div class="mt-2 mx-2">
          <Navbar infos={infos} />
        </div>
        <main class="flex flex-col container">{children}</main>
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
  );
}
