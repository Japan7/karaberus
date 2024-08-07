import { Command } from "@tauri-apps/api/shell";
import { HiSolidPlayCircle } from "solid-icons/hi";
import { Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { getSessionToken } from "../utils/session";

export default function TauriKaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  const downloadEndpoint = (type: string) =>
    `${location.origin}/api/kara/${props.kara.ID}/download/${type}`;

  const play: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    e.preventDefault();
    const args = [
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
    command.on("error", (error) => console.error(`command error: "${error}"`));
    command.stdout.on("data", (line) =>
      console.log(`command stdout: "${line}"`),
    );
    command.stderr.on("data", (line) =>
      console.log(`command stderr: "${line}"`),
    );
    const handle = await command.spawn();
    console.log(`mpv started with pid ${handle.pid}`);
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
