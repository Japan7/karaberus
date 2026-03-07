import { A } from "@solidjs/router";
import { FaSolidInfinity } from "solid-icons/fa";
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

function revision_url() {
  return `${import.meta.env.VITE_REPOSITORY_URL}/commit/${import.meta.env.VITE_REVISION}`;
}

export default function Sidebar(props: { closeDrawer: () => void }) {
  return (
    <div class="flex menu menu-lg bg-base-200 text-base-content min-h-full w-80 p-4">
      <ul class="flex-none">
        <li>
          <h2 class="menu-title flex gap-x-2">
            <HiOutlineMusicalNote class="size-5" />
            Karaokes
          </h2>
          <ul class="space-y-1">
            <li onclick={() => props.closeDrawer()}>
              <A href="/karaoke/new" activeClass="active">
                <HiSolidPencilSquare class="size-5" />
                New Karaoke
              </A>
            </li>
            <li onclick={() => props.closeDrawer()}>
              <A href="/karaoke/mugen" activeClass="active">
                <FaSolidInfinity class="size-5" />
                Mugen Import
              </A>
            </li>
            <li onclick={() => props.closeDrawer()}>
              <A href="/karaoke/browse" activeClass="active">
                <HiSolidGlobeAsiaAustralia class="size-5" />
                Browse
              </A>
            </li>
            <li onclick={() => props.closeDrawer()} class="disabled">
              {/* <A href="/karaoke/issues" activeClass="active"> */}
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
            <li onclick={() => props.closeDrawer()}>
              <A href="/tags/media" activeClass="active">
                <HiSolidTv class="size-5" />
                Media
              </A>
            </li>
            <li onclick={() => props.closeDrawer()}>
              <A href="/tags/artist" activeClass="active">
                <HiSolidMicrophone class="size-5" />
                Artist
              </A>
            </li>
            <li onclick={() => props.closeDrawer()}>
              <A href="/tags/author" activeClass="active">
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
            <li onclick={() => props.closeDrawer()}>
              <A href="/fonts/browse" activeClass="active">
                <HiSolidGlobeAsiaAustralia class="size-5" />
                Browse
              </A>
            </li>
          </ul>
        </li>
      </ul>
      <div class="flex-auto"></div>
      <div class="flex-none">
        <p>
          Karaberus{" "}
          <a class="link" href={revision_url()}>
            {import.meta.env.VITE_REVISION.substring(0, 8)}
          </a>
        </p>
      </div>
    </div>
  );
}
