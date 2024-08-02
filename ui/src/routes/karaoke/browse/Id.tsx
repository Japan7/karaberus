import { useParams } from "@solidjs/router";
import { onMount } from "solid-js";
import { fileForm, karaberus } from "../../../utils/karaberus-client";

export default function KaraokeBrowseId() {
  const params = useParams();

  let playerRef!: HTMLVideoElement;

  const upload = async (filetype: string, file?: File) => {
    if (!file) return;

    const resp = await karaberus.PUT("/api/kara/{id}/upload/{filetype}", {
      params: {
        path: {
          id: parseInt(params.id),
          filetype,
        },
      },
      ...fileForm(file),
    });

    if (resp.error) {
      alert(resp.error);
      return;
    }

    alert(filetype + " uploaded successfully");
    location.reload();
  };

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
      subUrl: `/api/kara/${params.id}/download/sub`,
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
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    window.octopusInstance = new SubtitlesOctopus(options);
  };

  return (
    <>
      <input
        type="file"
        onchange={(e) => upload("video", e.target.files?.[0])}
        class="file-input file-input-bordered w-full max-w-xs"
      />
      <input
        type="file"
        onchange={(e) => upload("sub", e.target.files?.[0])}
        class="file-input file-input-bordered w-full max-w-xs"
      />

      <video
        src={`/api/kara/${params.id}/download/video`}
        controls
        loop
        ref={playerRef}
      />
    </>
  );
}
