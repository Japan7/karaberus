import { invoke } from "@tauri-apps/api/core";
import { HiSolidPlayCircle } from "solid-icons/hi";
import { createSignal, Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { apiUrl } from "../utils/karaberus-client";
import { getPlayerToken } from "../utils/session";

export default function MpvKaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  const [getLoading, setLoading] = createSignal(false);

  const downloadEndpoint = (type: string) =>
    apiUrl(`api/kara/${props.kara.ID}/download/${type}`);

  const play: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    e.preventDefault();
    setLoading(true);
    const token = await getPlayerToken();
    await invoke("play_mpv", {
      video: props.kara.VideoUploaded ? downloadEndpoint("video") : undefined,
      inst: props.kara.InstrumentalUploaded
        ? downloadEndpoint("inst")
        : undefined,
      sub: props.kara.SubtitlesUploaded ? downloadEndpoint("sub") : undefined,
      token,
    });
    setLoading(false);
  };

  return (
    <Show when={props.kara.VideoUploaded || props.kara.InstrumentalUploaded}>
      <button disabled={getLoading()} onclick={play} class="btn btn-primary">
        <Show
          when={!getLoading()}
          fallback={<span class="loading loading-spinner loading-lg" />}
        >
          <HiSolidPlayCircle class="size-5" />
          Add to mpv queue
        </Show>
      </button>
    </Show>
  );
}
