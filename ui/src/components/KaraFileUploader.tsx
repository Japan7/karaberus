import { HiOutlineArrowTopRightOnSquare } from "solid-icons/hi";
import { createSignal, Show } from "solid-js";

export default function KaraFileUploader(props: {
  title: string;
  putUrl: string;
  downloadUrl?: string;
}) {
  const [getProgress, setProgress] = createSignal(0);

  const upload = async (file: File | undefined) => {
    if (!file) return;

    // xhr strikes again
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", props.putUrl);

    xhr.upload.addEventListener("progress", (event) => {
      setProgress(event.loaded / event.total);
    });

    xhr.addEventListener("load", () => {
      if (xhr.status === 200) {
        alert("Upload complete!");
        location.reload();
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
          <Show when={props.downloadUrl}>
            {(getDownloadUrl) => (
              <span class="label-text-alt">
                <a
                  href={getDownloadUrl()}
                  target="_blank"
                  rel="noreferrer"
                  class="link flex gap-x-1"
                >
                  Open current <HiOutlineArrowTopRightOnSquare class="size-4" />
                </a>
              </span>
            )}
          </Show>
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
