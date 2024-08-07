import { invoke } from "@tauri-apps/api";
import { HiSolidPlayCircle } from "solid-icons/hi";
import { Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { getSessionToken } from "../utils/session";

export default function TauriKaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  const downloadEndpoint = (type: string) =>
    `${location.origin}/api/kara/${props.kara.ID}/download/${type}`;
  const getVideoSrc = () =>
    props.kara.VideoUploaded ? downloadEndpoint("video") : undefined;
  const getInstSrc = () =>
    props.kara.InstrumentalUploaded ? downloadEndpoint("inst") : undefined;
  const getSubSrc = () =>
    props.kara.SubtitlesUploaded ? downloadEndpoint("sub") : undefined;

  const play: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    e.preventDefault();
    await invoke("play", {
      auth: getSessionToken(),
      video: getVideoSrc(),
      inst: getInstSrc(),
      sub: getSubSrc(),
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
