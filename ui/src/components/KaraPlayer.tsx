import { onMount } from "solid-js";

export default function KaraPlayer({ id }: { id: number | string }) {
  let playerRef!: HTMLVideoElement;

  onMount(() => {
    const subtitlesOctopusScript = document.createElement("script");
    subtitlesOctopusScript.src = "/libass-wasm/subtitles-octopus.js";
    subtitlesOctopusScript.async = true;
    subtitlesOctopusScript.onload = subtitlesOctopusInstantiate;
    document.head.appendChild(subtitlesOctopusScript);
  });

  const subtitlesOctopusInstantiate = () => {
    const options = {
      video: playerRef,
      subUrl: `/api/kara/${id}/download/sub`,
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
    // @ts-expect-error: global variable
    window.octopusInstance = new SubtitlesOctopus(options);
  };

  return (
    <video
      src={`/api/kara/${id}/download/video`}
      controls
      // @ts-expect-error: https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/controlsList
      controlslist="nofullscreen"
      playsinline
      loop
      ref={playerRef}
    />
  );
}
