import { A } from "@solidjs/router";
import { Icon } from "solid-heroicons";
import { atSymbol, musicalNote } from "solid-heroicons/outline";
import {
  arrowUpTray,
  exclamationCircle,
  globeAsiaAustralia,
} from "solid-heroicons/solid";
import routes from "../utils/routes";

export default function Sidebar() {
  return (
    <ul class="menu menu-lg bg-base-200 text-base-content min-h-full w-80 p-4">
      <li>
        <h2 class="menu-title flex gap-x-2">
          <Icon path={musicalNote} class="size-5" />
          Karaokes
        </h2>
        <ul class="space-y-1">
          <li>
            <A href={routes.KARAOKE_UPLOAD} activeClass="active">
              <Icon path={arrowUpTray} class="size-5" />
              Upload
            </A>
          </li>
          <li>
            <A href={routes.KARAOKE_BROWSE} activeClass="active">
              <Icon path={globeAsiaAustralia} class="size-5" />
              Browse
            </A>
          </li>
          <li class="disabled">
            <a>
              <Icon path={exclamationCircle} class="size-5" />
              Issues
            </a>
          </li>
        </ul>
      </li>
      <li>
        <h2 class="menu-title flex gap-x-2">
          <Icon path={atSymbol} class="size-5" />
          Fonts
        </h2>
        <ul class="space-y-1">
          <li>
            <A href={routes.FONTS_UPLOAD} activeClass="active">
              <Icon path={arrowUpTray} class="size-5" />
              Upload
            </A>
          </li>
          <li>
            <A href={routes.FONTS_BROWSE} activeClass="active">
              <Icon path={globeAsiaAustralia} class="size-5" />
              Browse
            </A>
          </li>
        </ul>
      </li>
    </ul>
  );
}
