import { invoke } from "@tauri-apps/api";
import { Show } from "solid-js";
import type { components } from "../utils/karaberus";
import { getSessionToken } from "../utils/session";

export default function TauriKaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  const getVideoSrc = () => `/api/kara/${props.kara.ID}/download/video`;
  const getSubSrc = () => `/api/kara/${props.kara.ID}/download/sub`;

  const play = async () => {
    await invoke("play", {
      video: location.origin + getVideoSrc(),
      sub: location.origin + getSubSrc(),
      auth: getSessionToken(),
    });
  };

  return (
    <Show when={props.kara.VideoUploaded}>
      <button onclick={play} class="btn btn-wide">
        Open in mpv
      </button>
    </Show>
  );
}
