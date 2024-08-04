import { HiOutlineBars3 } from "solid-icons/hi";
import type { KaraberusJwtPayload } from "../utils/session";
import ProfileDropdown from "./ProfileDropdown";

export default function Navbar({ infos }: { infos: KaraberusJwtPayload }) {
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
      <div class="flex-none">
        <ProfileDropdown infos={infos} />
      </div>
    </div>
  );
}
