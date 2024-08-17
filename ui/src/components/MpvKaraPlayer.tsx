import { invoke } from "@tauri-apps/api/core";
import { HiSolidPlayCircle } from "solid-icons/hi";
import { Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { apiUrl } from "../utils/karaberus-client";
import { getSessionToken } from "../utils/session";

export default function MpvKaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  const downloadEndpoint = (type: string) =>
    apiUrl(`api/kara/${props.kara.ID}/download/${type}`);

  const play: JSX.EventHandler<HTMLButtonElement, MouseEvent> = (e) => {
    e.preventDefault();
    invoke("play_mpv", {
      auth: getSessionToken(),
      video: props.kara.VideoUploaded ? downloadEndpoint("video") : undefined,
      inst: props.kara.InstrumentalUploaded
        ? downloadEndpoint("inst")
        : undefined,
      sub: props.kara.SubtitlesUploaded ? downloadEndpoint("sub") : undefined,
    });
  };

  return (
    <Show when={props.kara.VideoUploaded || props.kara.InstrumentalUploaded}>
      <button onclick={play} class="btn btn-lg btn-primary">
        <HiSolidPlayCircle class="size-5" />
        Open in mpv
      </button>
    </Show>
  );
}
