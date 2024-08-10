import { invoke } from "@tauri-apps/api/core";
import { HiSolidPlayCircle } from "solid-icons/hi";
import { createSignal, Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { apiUrl } from "../utils/karaberus-client";
import { getSessionToken } from "../utils/session";

export default function MpvKaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  const [getLoading] = createSignal(false);
  const [getStarted] = createSignal(false);

  const downloadEndpoint = (type: string) =>
    apiUrl(`api/kara/${props.kara.ID}/download/${type}`);

  const play: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
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
      <button
        disabled={getLoading() || getStarted()}
        onclick={play}
        class="btn btn-lg btn-primary"
      >
        <Show
          when={!getLoading()}
          fallback={<span class="loading loading-spinner loading-lg" />}
        >
          <HiSolidPlayCircle class="size-5" />
          Open in mpv
        </Show>
      </button>
    </Show>
  );
}
