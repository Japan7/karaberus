import { HiOutlineArrowDownTray, HiOutlineBars3 } from "solid-icons/hi";
import { Show } from "solid-js";
import type { KaraberusJwtPayload } from "../utils/session";
import { isTauri, RELEASE_URL } from "../utils/tauri";
import ProfileDropdown from "./ProfileDropdown";

export default function Navbar(props: { infos: KaraberusJwtPayload }) {
  return (
    <div class="navbar bg-base-100 shadow-xl rounded-box">
      <div class="flex-none">
        <label for="main-drawer" class="btn btn-square btn-ghost lg:hidden">
          <HiOutlineBars3 class="size-6" />
        </label>
      </div>
      <div class="flex-1">
        <a href="/" class="btn btn-ghost text-xl">
          Karaberus
        </a>
      </div>
      <div class="flex-none gap-x-2">
        <Show when={!isTauri()}>
          <a
            href={RELEASE_URL}
            target="_blank"
            rel="noopener noreferrer"
            class="hidden sm:btn"
          >
            <HiOutlineArrowDownTray class="size-5" />
            Desktop app
          </a>
        </Show>
        <ProfileDropdown infos={props.infos} />
      </div>
    </div>
  );
}
