import { createSignal, Show, type JSX } from "solid-js";

const CHUNK_SIZE = 5 * 1024 * 1024;

export default function FileUploader(props: {
  title?: string;
  method: string;
  url: string;
  chunked?: boolean;
  onUpload: () => void;
  altChildren?: JSX.Element;
}) {
  const [getProgress, setProgress] = createSignal(0);

  const upload = (file: File | undefined) => {
    if (!file) return;

    const formData = new FormData();
    formData.append("file", file);

    // xhr strikes again
    const xhr = new XMLHttpRequest();

    xhr.open(props.method, props.url);

    xhr.upload.addEventListener("progress", (event) => {
      setProgress(event.loaded / event.total);
    });

    xhr.addEventListener("load", () => {
      setProgress(0);
      if (xhr.status === 200) {
        props.onUpload();
      } else {
        alert(xhr.responseText);
      }
    });

    setProgress(0);
    xhr.send(formData);
  };

  const chunkedUpload = (file: File | undefined) => {
    if (!file) return;

    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    let currentChunk = 0;

    const uploadNextChunk = () => {
      const start = currentChunk * CHUNK_SIZE;
      const end = Math.min(start + CHUNK_SIZE, file.size);
      const chunk = file.slice(start, end);

      const formData = new FormData();
      formData.append("file", chunk, file.name);

      // xhr strikes again
      const xhr = new XMLHttpRequest();

      xhr.open(props.method, props.url);

      xhr.setRequestHeader(
        "Content-Range",
        `bytes ${start}-${end - 1}/${file.size}`,
      );

      xhr.upload.addEventListener("progress", (e) => {
        setProgress(start + e.loaded / file.size);
      });

      xhr.addEventListener("load", () => {
        if (xhr.status === 200) {
          if (++currentChunk < totalChunks) {
            uploadNextChunk();
          } else {
            setProgress(0);
            props.onUpload();
          }
        } else {
          alert(xhr.responseText);
        }
      });

      xhr.send(formData);
    };

    setProgress(0);
    uploadNextChunk();
  };

  return (
    <div class="flex flex-col gap-y-2">
      <label class="form-control">
        <div class="label">
          <span class="label-text">{props.title}</span>
          <span class="label-text-alt">{props.altChildren}</span>
        </div>
        <input
          type="file"
          onchange={(e) =>
            (props.chunked ? chunkedUpload : upload)(e.target.files?.[0])
          }
          class="file-input file-input-bordered"
        />
      </label>

      <Show when={getProgress()}>
        <div class="flex items-center gap-x-1">
          <progress value={getProgress()} class="progress" />
          <pre>{(getProgress() * 100).toFixed(1).padStart(5)}%</pre>
        </div>
      </Show>
    </div>
  );
}
