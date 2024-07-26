import type { JSX } from "solid-js";
import routes from "../utils/routes";
import { clearSession, type KaraberusTokenPayload } from "../utils/session";

export default function ProfileDropdown({
  infos,
}: {
  infos: KaraberusTokenPayload;
}) {
  const logout: JSX.EventHandler<HTMLAnchorElement, MouseEvent> = () => {
    clearSession();
    location.href = routes.HOME;
  };

  return (
    <div class="dropdown dropdown-bottom dropdown-end">
      <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
        <div class="w-10 rounded-full">
          <img alt={infos.name} src={infos.picture} />
        </div>
      </div>
      <ul
        tabindex="0"
        class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow"
      >
        <li class="disabled">
          <a>Profile</a>
        </li>
        <li class="disabled">
          <a>Settings</a>
        </li>
        <li>
          <a onclick={logout}>Logout</a>
        </li>
      </ul>
    </div>
  );
}
