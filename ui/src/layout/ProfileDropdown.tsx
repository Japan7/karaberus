import {
  clearSession,
  isAdmin,
  type KaraberusJwtPayload,
} from "../utils/session";

export default function ProfileDropdown(props: { infos: KaraberusJwtPayload }) {
  const logout = () => {
    clearSession();
    location.reload();
  };

  const closeDropdown = () => {
    const elem = document.activeElement;
    if (elem instanceof HTMLElement) {
      elem.blur();
    }
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
        <li onclick={() => closeDropdown()}>
          <a href="/profile" class="flex flex-col gap-0">
            <span class="w-full font-bold text-lg">{props.infos.name}</span>
            <span class="w-full font-bold">{isAdmin() ? "Admin" : "User"}</span>
          </a>
        </li>
        <li onclick={() => closeDropdown()}>
          <a href="/settings">Settings</a>
        </li>
        <li onclick={() => closeDropdown()}>
          <a onclick={() => logout()}>Logout</a>
        </li>
      </ul>
    </div>
  );
}
