import SubtitlesOctopus from "libass-wasm";
import { createEffect, onCleanup } from "solid-js";

export default function KaraPlayer(props: { id: number | string }) {
  let playerRef!: HTMLVideoElement;

  let octopus: SubtitlesOctopus | undefined;

  const videoSrc = `/api/kara/${props.id}/download/video`;
  const subSrc = `/api/kara/${props.id}/download/sub`;

  createEffect(() => {
    onCleanup(() => octopus?.dispose());
  });

  const setupOctopus = () => {
    octopus?.dispose();
    const options = {
      video: playerRef,
      subUrl: subSrc,
      fonts: [
        "/amaranth/amaranth-latin-400-italic.woff",
        "/amaranth/amaranth-latin-400-italic.woff2",
        "/amaranth/amaranth-latin-400-normal.woff",
        "/amaranth/amaranth-latin-400-normal.woff2",
        "/amaranth/amaranth-latin-700-italic.woff",
        "/amaranth/amaranth-latin-700-italic.woff2",
        "/amaranth/amaranth-latin-700-normal.woff",
        "/amaranth/amaranth-latin-700-normal.woff2",
      ],
      workerUrl: "/libass-wasm/subtitles-octopus-worker.js",
      legacyWorkerUrl: "/libass-wasm/subtitles-octopus-worker-legacy.js",
    };
    octopus = new SubtitlesOctopus(options);
  };

  return (
    <video
      src={videoSrc}
      controls
      // @ts-expect-error: https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/controlsList
      controlslist="nofullscreen"
      playsinline
      loop
      oncanplay={setupOctopus}
      onerror={async () => {
        await new Promise((resolve) => setTimeout(resolve, 1000));
        playerRef.src = videoSrc;
      }}
      ref={playerRef}
      class="rounded-2xl"
    />
  );
}
