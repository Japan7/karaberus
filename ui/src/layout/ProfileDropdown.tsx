import type { JSX } from "solid-js";
import {
  clearSession,
  isAdmin,
  type KaraberusJwtPayload,
} from "../utils/session";

export default function ProfileDropdown(props: { infos: KaraberusJwtPayload }) {
  let detailsRef!: HTMLDetailsElement;

  const logout: JSX.EventHandler<HTMLElement, MouseEvent> = () => {
    clearSession();
    location.href = "/";
  };

  return (
    <details ref={detailsRef} class="dropdown dropdown-bottom dropdown-end">
      <summary class="btn btn-ghost btn-circle avatar">
        <div class="w-10 rounded-full">
          <img alt={props.infos.name} src={props.infos.picture} />
        </div>
      </summary>
      <ul class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow">
        <li onclick={() => detailsRef.removeAttribute("open")}>
          <a href="/profile" class="flex flex-col gap-0">
            <span class="w-full font-bold text-lg">{props.infos.name}</span>
            <span class="w-full font-bold">{isAdmin() ? "Admin" : "User"}</span>
          </a>
        </li>
        <li onclick={() => detailsRef.removeAttribute("open")}>
          <a href="/settings">Settings</a>
        </li>
        <li onclick={() => detailsRef.removeAttribute("open")}>
          <a onclick={logout}>Logout</a>
        </li>
      </ul>
    </details>
  );
}
