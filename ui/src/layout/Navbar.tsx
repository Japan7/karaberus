import { Icon } from "solid-heroicons";
import { bars_3 } from "solid-heroicons/outline";
import routes from "../utils/routes";
import type { KaraberusTokenPayload } from "../utils/session";
import ProfileDropdown from "./ProfileDropdown";

export default function Navbar({ infos }: { infos: KaraberusTokenPayload }) {
  return (
    <div class="navbar bg-base-100 shadow-xl rounded-box">
      <div class="flex-none">
        <label for="main-drawer" class="btn btn-square btn-ghost lg:hidden">
          <Icon path={bars_3} class="size-6" />
        </label>
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
  );
}
