import { isTauri } from "@tauri-apps/api/core";
import { open } from "@tauri-apps/plugin-shell";
import type { JSX } from "solid-js";
import { buildRedirectUrl } from "../utils/tauri";

export default function DownloadAnchor(props: {
  href: string;
  download: string;
  children: JSX.Element;
}) {
  const onClick: JSX.EventHandler<HTMLAnchorElement, MouseEvent> = (e) => {
    e.preventDefault();
    open(buildRedirectUrl(props.href).toString());
  };

  return (
    <a
      href={isTauri() ? undefined : props.href}
      download={props.download}
      onclick={isTauri() ? onClick : undefined}
      rel="noreferrer"
      class="link flex gap-x-1"
    >
      {props.children}
    </a>
  );
}
