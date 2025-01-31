import { createSignal, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";

export default function AuthorEditor(props: {
  author?: components["schemas"]["TimingAuthor"];
  onSubmit: (author: components["schemas"]["AuthorInfo"]) => void;
  reset?: boolean;
}) {
  const [getName, setName] = createSignal(props.author?.Name ?? "");
  const [getPublicName, setPublicName] = createSignal(
    props.author?.PublicName ?? "",
  );

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = (e) => {
    e.preventDefault();
    props.onSubmit({
      name: getName(),
      external_name: getPublicName(),
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
          placeholder="bebou69"
          value={getName()}
          onInput={(e) => setName(e.currentTarget.value)}
          class="input input-bordered w-full"
        />
        <div class="label">
          <span class="label-text">External Name</span>
          <span class="label-text-alt">(optional)</span>
        </div>
        <input
          type="text"
          required
          placeholder="bebou69"
          value={getPublicName()}
          onInput={(e) => setPublicName(e.currentTarget.value)}
          class="input input-bordered w-full"
        />
      </label>

      <input type="submit" class="btn btn-primary" />
    </form>
  );
}
