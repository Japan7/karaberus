import { createSignal, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";

export default function ArtistEditor(props: {
  artist?: components["schemas"]["Artist"];
  onSubmit: (artist: components["schemas"]["ArtistInfo"]) => void;
  reset?: boolean;
}) {
  const [getName, setName] = createSignal(props.artist?.Name ?? "");
  const [getAdditionalNames, setAdditionalNames] = createSignal(
    props.artist?.AdditionalNames?.join("\n") ?? "",
  );

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = (e) => {
    e.preventDefault();
    props.onSubmit({
      name: getName(),
      additional_names: getAdditionalNames().trim().split("\n"),
    });
    if (props.reset) {
      (e.target as HTMLFormElement).reset();
    }
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
