import { Command } from "@tauri-apps/api/shell";
import { HiSolidPlayCircle } from "solid-icons/hi";
import { createSignal, Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { getSessionToken } from "../utils/session";

export default function MpvKaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  const [getLoading, setLoading] = createSignal(false);

  const downloadEndpoint = (type: string) =>
    `${location.origin}/api/kara/${props.kara.ID}/download/${type}`;

  const play: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    e.preventDefault();
    setLoading(true);
    const args = [
      "--quiet",
      `--http-header-fields=Authorization: Bearer ${getSessionToken()}`,
    ];
    if (props.kara.SubtitlesUploaded) {
      args.push(`--sub-file=${downloadEndpoint("sub")}`);
    }
    if (props.kara.VideoUploaded) {
      if (props.kara.InstrumentalUploaded) {
        args.push(`--external-file=${downloadEndpoint("inst")}`);
      }
      args.push(downloadEndpoint("video"));
    } else if (props.kara.InstrumentalUploaded) {
      args.push(downloadEndpoint("inst"));
    } else {
      return;
    }
    const command = new Command("mpv", args);
    command.on("close", (exit) => {
      console.debug(`command finished with code ${exit.code}`);
      if (exit.code !== 0) {
        alert("An error occurred");
      }
      setLoading(false);
    });
    command.stdout.on("data", console.log);
    command.stderr.on("data", console.error);
    const handle = await command.spawn();
    console.debug(`mpv started with pid ${handle.pid}`);
  };

  return (
    <Show when={props.kara.VideoUploaded || props.kara.InstrumentalUploaded}>
      <button
        disabled={getLoading()}
        onclick={play}
        class="btn btn-lg btn-primary"
      >
        <HiSolidPlayCircle class="size-5" />
        Open in mpv
      </button>
    </Show>
  );
}
