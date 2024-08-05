import { createSignal, Show, type JSX } from "solid-js";

export default function FileUploader(props: {
  title?: string;
  method: string;
  url: string;
  onUpload: () => void;
  altChildren?: JSX.Element;
}) {
  const [getProgress, setProgress] = createSignal(0);

  const upload = async (file: File | undefined) => {
    if (!file) return;

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

    const formData = new FormData();
    formData.append("file", file);
    xhr.send(formData);
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
          onchange={(e) => upload(e.target.files?.[0])}
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
