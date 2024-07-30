import { A } from "@solidjs/router";
import {
  HiOutlineAtSymbol,
  HiOutlineMusicalNote,
  HiOutlineTag,
  HiSolidExclamationCircle,
  HiSolidGlobeAsiaAustralia,
  HiSolidMicrophone,
  HiSolidPencilSquare,
  HiSolidTv,
  HiSolidUser,
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
            <A href={routes.KARAOKE_CREATE} activeClass="active">
              <HiSolidPencilSquare class="size-5" />
              Create
            </A>
          </li>
          <li>
            <A href={routes.KARAOKE_BROWSE} activeClass="active">
              <HiSolidGlobeAsiaAustralia class="size-5" />
              Browse
            </A>
          </li>
          <li class="disabled">
            {/* <A href={routes.KARAOKE_ISSUES} activeClass="active"> */}
            <a>
              <HiSolidExclamationCircle class="size-5" />
              Issues
            </a>
            {/* </A> */}
          </li>
        </ul>
      </li>
      <li>
        <h2 class="menu-title flex gap-x-2">
          <HiOutlineTag class="size-5" />
          Tags
        </h2>
        <ul class="space-y-1">
          <li>
            <A href={routes.TAGS_MEDIA} activeClass="active">
              <HiSolidTv class="size-5" />
              Media
            </A>
          </li>
          <li>
            <A href={routes.TAGS_ARTIST} activeClass="active">
              <HiSolidMicrophone class="size-5" />
              Artist
            </A>
          </li>
          <li>
            <A href={routes.TAGS_AUTHOR} activeClass="active">
              <HiSolidUser class="size-5" />
              Author
            </A>
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
