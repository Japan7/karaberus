import { A } from "@solidjs/router";
import {
  HiOutlineAtSymbol,
  HiOutlineMusicalNote,
  HiSolidArrowUpTray,
  HiSolidExclamationCircle,
  HiSolidGlobeAsiaAustralia,
} from "solid-icons/hi";
import routes from "../utils/routes";

export default function Sidebar() {
  return (
    <ul class="menu menu-lg bg-base-200 text-base-content min-h-full w-80 p-4">
      <li>
        <h2 class="menu-title flex gap-x-2">
          <HiOutlineMusicalNote class="size-5" />
          Karaokes
        </h2>
        <ul class="space-y-1">
          <li>
            <A href={routes.KARAOKE_UPLOAD} activeClass="active">
              <HiSolidArrowUpTray class="size-5" />
              Upload
            </A>
          </li>
          <li>
            <A href={routes.KARAOKE_BROWSE} activeClass="active">
              <HiSolidGlobeAsiaAustralia class="size-5" />
              Browse
            </A>
          </li>
          <li class="disabled">
            <a>
              <HiSolidExclamationCircle class="size-5" />
              Issues
            </a>
          </li>
        </ul>
      </li>
      <li>
        <h2 class="menu-title flex gap-x-2">
          <HiOutlineAtSymbol class="size-5" />
          Fonts
        </h2>
        <ul class="space-y-1">
          <li>
            <A href={routes.FONTS_UPLOAD} activeClass="active">
              <HiSolidArrowUpTray class="size-5" />
              Upload
            </A>
          </li>
          <li>
            <A href={routes.FONTS_BROWSE} activeClass="active">
              <HiSolidGlobeAsiaAustralia class="size-5" />
              Browse
            </A>
          </li>
        </ul>
      </li>
    </ul>
  );
}
