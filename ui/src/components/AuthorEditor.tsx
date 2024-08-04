import { createSignal, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";

export default function AuthorEditor(props: {
  onAdd: (author: components["schemas"]["TimingAuthor"]) => void;
}) {
  const [getName, setName] = createSignal("");

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();
    const resp = await karaberus.POST("/api/tags/author", {
      body: { name: getName() },
    });
    if (resp.error) {
      alert(resp.error);
      return;
    }
    (e.target as HTMLFormElement).reset();
    props.onAdd(resp.data.author);
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full max-w-xs">
      <label>
        <div class="label">
          <span class="label-text">Name</span>
          <span class="label-text-alt">(required)</span>
        </div>
        <input
          type="text"
          required
          value={getName()}
          onInput={(e) => setName(e.currentTarget.value)}
          class="input input-bordered w-full"
        />
      </label>

      <input type="submit" class="btn" />
    </form>
  );
}
