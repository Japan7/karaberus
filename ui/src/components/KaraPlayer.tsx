import wasmURL from "@ffmpeg/core-mt/wasm?url";
import coreURL from "@ffmpeg/core-mt?url";
import { FFmpeg } from "@ffmpeg/ffmpeg";
import { fetchFile } from "@ffmpeg/util";
import { createSignal, Match, onMount, Show, Switch } from "solid-js";

export default function KaraPlayer({ id }: { id: number | string }) {
  const [getFFmpeg, setFFmpeg] = createSignal<FFmpeg>();
  const [getProgress, setProgress] = createSignal<number>();
  const [getSrc, setSrc] = createSignal<string>();

  onMount(async () => {
    const ffmpeg = new FFmpeg();
    ffmpeg.on("log", ({ message }) => {
      console.log(message);
    });
    ffmpeg.on("progress", ({ progress }) => {
      setProgress(progress);
    });
    await ffmpeg.load({
      coreURL,
      wasmURL,
      workerURL: "/ffmpeg-core.worker.js",
    });
    setFFmpeg(ffmpeg);
    console.log("FFmpeg is ready");
  });

  const transcode = async () => {
    const ffmpeg = getFFmpeg();
    if (!ffmpeg) {
      console.error("FFmpeg is not ready");
      return;
    }
    const video = await fetchFile(`/api/kara/${id}/download/video`);
    const sub = await fetchFile(`/api/kara/${id}/download/sub`);
    const font = await fetchFile("/Amaranth-Regular.ttf");
    await ffmpeg.writeFile("video", video);
    await ffmpeg.writeFile("sub", sub);
    await ffmpeg.writeFile("/tmp/Amaranth-Regular.ttf", font);
    ffmpeg.exec([
      "-i",
      "video",
      "-vf",
      "subtitles=sub:fontsdir=/tmp:force_style='Fontname='Amaranth'",
      "-preset",
      "ultrafast",
      "output.mp4",
    ]);
    const data = await ffmpeg.readFile("output.mp4");
    setSrc(URL.createObjectURL(new Blob([data], { type: "video/mp4" })));
  };

  return (
    <Switch
      fallback={
        <Show when={getFFmpeg()} fallback={<p>Loading FFmpeg...</p>}>
          <button onClick={transcode} class="btn">
            Preview
          </button>
        </Show>
      }
    >
      <Match when={getSrc()}>
        {(getSrc) => (
          <video src={getSrc()} controls autoplay loop playsinline />
        )}
      </Match>
      <Match when={getProgress()}>
        {(getProgress) => (
          <div class="flex items-center gap-x-2">
            <progress value={getProgress()} class="progress" />
            <p>{Math.round(getProgress() * 100)}%</p>
          </div>
        )}
      </Match>
    </Switch>
  );
}
