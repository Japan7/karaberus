import { createSignal, Show, type JSX } from "solid-js";
import { getSessionToken } from "../utils/session";
import { IS_TAURI_DIST_BUILD } from "../utils/tauri";

export default function FileUploader(props: {
  title?: string;
  method: string;
  url: string;
  onUpload: () => void;
  altChildren?: JSX.Element;
}) {
  const [getProgress, setProgress] = createSignal(0);

  const upload: JSX.EventHandler<HTMLInputElement, Event> = (e) => {
    const target = e.target as HTMLInputElement;

    const file = target.files?.[0];
    if (!file) return;

    // xhr strikes again
    const xhr = new XMLHttpRequest();
    xhr.open(props.method, props.url);
    if (IS_TAURI_DIST_BUILD) {
      xhr.setRequestHeader("Authorization", `JWT ${getSessionToken()}`);
    }

    xhr.upload.addEventListener("progress", (event) => {
      setProgress(event.loaded / event.total);
    });

    xhr.addEventListener("load", () => {
      target.value = "";
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
    <label class="form-control">
      <div class="label grid grid-cols-2">
        <span class="label-text">{props.title}</span>
        <span class="label-text-alt">
          <Show when={getProgress()}>
            <div class="flex items-center gap-x-1">
              <progress
                value={getProgress()}
                class="progress progress-success"
              />
              <pre>{(getProgress() * 100).toFixed(1).padStart(5)}%</pre>
            </div>
          </Show>
        </span>
      </div>
      <input
        type="file"
        onchange={upload}
        class="file-input file-input-bordered"
      />
      <div class="label">
        <span class="label-text-alt">{props.altChildren}</span>
      </div>
    </label>
  );
}
