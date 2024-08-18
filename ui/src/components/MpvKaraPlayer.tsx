import { invoke } from "@tauri-apps/api/core";
import { HiOutlineCheck, HiSolidPlayCircle } from "solid-icons/hi";
import { createSignal, Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { apiUrl, karaberus } from "../utils/karaberus-client";
import { getStore } from "../utils/tauri";

export default function MpvKaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  const [getLoading, setLoading] = createSignal(false);
  const [getSuccess, setSuccess] = createSignal(false);

  const downloadEndpoint = (type: string) =>
    apiUrl(`api/kara/${props.kara.ID}/download/${type}`);

  const play: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    e.preventDefault();
    setLoading(true);
    await ensurePlayerToken();
    await invoke("play_mpv", {
      video: props.kara.VideoUploaded ? downloadEndpoint("video") : undefined,
      inst: props.kara.InstrumentalUploaded
        ? downloadEndpoint("inst")
        : undefined,
      sub: props.kara.SubtitlesUploaded ? downloadEndpoint("sub") : undefined,
      title: props.kara.Title,
    });
    setLoading(false);
    setSuccess(true);
    setTimeout(() => setSuccess(false), 1000);
  };

  const ensurePlayerToken = async () => {
    const store = getStore();
    let token = await store.get<string>("player_token");
    if (!token) {
      const resp = await karaberus.POST("/api/token", {
        body: {
          name: "karaberus_player",
          scopes: { kara: false, kara_ro: true, user: false },
        },
      });
      if (resp.error) {
        throw new Error(resp.error.title);
      }
      token = resp.data.token;
      await store.set("player_token", token);
      await store.save();
    }
    return token;
  };

  return (
    <Show when={props.kara.VideoUploaded || props.kara.InstrumentalUploaded}>
      <button
        disabled={getLoading()}
        onclick={play}
        class="btn btn-primary"
        classList={{ "btn-success": getSuccess() }}
      >
        <Show
          when={!getLoading()}
          fallback={<span class="loading loading-spinner" />}
        >
          {getSuccess() ? (
            <HiOutlineCheck class="size-5" />
          ) : (
            <HiSolidPlayCircle class="size-5" />
          )}
          Add to mpv queue
        </Show>
      </button>
    </Show>
  );
}
