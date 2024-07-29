import { createSignal, type JSX } from "solid-js";
import { fileForm, karaberus } from "../../utils/karaberus-client";

export default function FontsUpload() {
  const [getFile, setFile] = createSignal<File>();

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();

    const file = getFile();
    if (!file) {
      alert("Please select a file");
      return;
    }

    const resp = await karaberus.POST("/api/font", fileForm(file));
    if (resp.error) {
      alert(resp.error);
    } else {
      alert("Font uploaded");
    }
  };

  return (
    <>
      <h1 class="header">Upload Font</h1>

      <form onsubmit={onsubmit} class="flex gap-x-2">
        <input
          type="file"
          required
          onchange={(e) => setFile(e.target.files?.[0])}
          class="file-input file-input-bordered w-full max-w-xs"
        />
        <input type="submit" class="btn" />
      </form>
    </>
  );
}
