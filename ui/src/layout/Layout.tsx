import { Icon } from "solid-heroicons";
import { bars_3 } from "solid-heroicons/outline";
import type { JSX } from "solid-js";
import routes from "../utils/routes";
import { getSessionInfos } from "../utils/session";
import ProfileDropdown from "./ProfileDropdown";
import AuthHero from "./AuthHero";

export default function Layout({ children }: { children: JSX.Element }) {
  const infos = getSessionInfos();

  return infos ? (
    <>
      <div class="m-2">
        <div class="navbar bg-base-100 shadow-xl rounded-box">
          <div class="flex-none">
            <button class="btn btn-square btn-ghost">
              <Icon path={bars_3} class="size-6" />
            </button>
          </div>
          <div class="flex-1">
            <a href={routes.HOME} class="btn btn-ghost text-xl">
              Karaberus
            </a>
          </div>
          <div class="flex-none">
            <ProfileDropdown infos={infos} />
          </div>
        </div>
      </div>
      {children}
    </>
  ) : (
    <AuthHero />
  );
}
