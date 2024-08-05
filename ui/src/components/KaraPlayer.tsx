import SubtitlesOctopus from "libass-wasm";
import { createEffect, onCleanup, Show } from "solid-js";
import type { components } from "../utils/karaberus";

export default function KaraPlayer(props: {
  kara: components["schemas"]["KaraInfoDB"];
}) {
  let playerRef!: HTMLVideoElement;

  let octopus: SubtitlesOctopus | undefined;

  const getVideoSrc = () => `/api/kara/${props.kara.ID}/download/video`;
  const getSubSrc = () => `/api/kara/${props.kara.ID}/download/sub`;

  createEffect(() => {
    onCleanup(() => {
      octopus?.dispose();
      octopus = undefined;
    });
  });

  const setupOctopus = () => {
    if (props.kara.SubtitlesUploaded && !octopus) {
      octopus = new SubtitlesOctopus({
        video: playerRef,
        subUrl: getSubSrc(),
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
      });
    }
  };

  return (
    <Show when={props.kara.VideoUploaded}>
      <video
        src={getVideoSrc()}
        controls
        // @ts-expect-error: https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/controlsList
        controlslist="nofullscreen"
        playsinline
        autoplay
        loop
        oncanplay={setupOctopus}
        ref={playerRef}
        class="rounded-2xl"
      />
    </Show>
  );
}
