import type { JSX } from "solid-js";
import { clearSession, type KaraberusJwtPayload } from "../utils/session";

export default function ProfileDropdown(props: { infos: KaraberusJwtPayload }) {
  const logout: JSX.EventHandler<HTMLElement, MouseEvent> = () => {
    clearSession();
    location.href = "/";
  };

  return (
    <div class="dropdown dropdown-bottom dropdown-end">
      <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
        <div class="w-10 rounded-full">
          <img alt={props.infos.name} src={props.infos.picture} />
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
