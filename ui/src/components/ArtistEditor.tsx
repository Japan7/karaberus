import { createSignal, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";

export default function ArtistEditor(props: {
  onAdd: (artist: components["schemas"]["Artist"]) => void;
}) {
  const [getName, setName] = createSignal("");
  const [getAdditionalNames, setAdditionalNames] = createSignal("");

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();
    const resp = await karaberus.POST("/api/tags/artist", {
      body: {
        name: getName(),
        additional_names: getAdditionalNames().trim().split("\n"),
      },
    });
    if (resp.error) {
      alert(resp.error.title);
      return;
    }
    (e.target as HTMLFormElement).reset();
    props.onAdd(resp.data.artist);
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full">
      <label>
        <div class="label">
          <span class="label-text">Name</span>
          <span class="label-text-alt">(required)</span>
        </div>
        <input
          type="text"
          required
          placeholder="Yoko Takahashi"
          value={getName()}
          onInput={(e) => setName(e.currentTarget.value)}
          class="input input-bordered w-full"
        />
      </label>

      <label>
        <div class="label">
          <span class="label-text">Additional names</span>
          <span class="label-text-alt">1 per line</span>
        </div>
        <textarea
          placeholder={"YAWMIN\n高橋洋子"}
          value={getAdditionalNames()}
          onInput={(e) => setAdditionalNames(e.currentTarget.value)}
          class="textarea textarea-bordered w-full"
        />
      </label>

      <input type="submit" class="btn btn-primary" />
    </form>
  );
}
