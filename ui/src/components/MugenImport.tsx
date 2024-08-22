import { createEffect, createSignal, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";

export default function MugenImport(props: {
  onImport: (kara: components["schemas"]["MugenImport"]) => void;
}) {
  const [getInput, setInput] = createSignal("");
  const [getName, setName] = createSignal("");

  const getKid = () => getInput().split("/").pop() ?? getInput();

  createEffect(async () => {
    if (getKid().length !== 36) {
      setName("");
      return;
    }
    const resp = await fetch(`https://kara.moe/api/karas/${getKid()}`);
    const json = await resp.json();
    setName(json.karafile);
  });

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();

    const resp = await karaberus.POST("/api/mugen", {
      body: {
        mugen_kid: getKid(),
      },
    });

    if (resp.error) {
      alert(resp.error.title);
      return;
    }

    (e.target as HTMLFormElement).reset();
    props.onImport(resp.data.import);
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full">
      <label>
        <div class="label">
          <span class="label-text">Karaoke ID/URL</span>
          <span class="label-text-alt">(required)</span>
        </div>
        <input
          type="text"
          required
          placeholder="https://kara.moe/kara/zankoku-na-tenshi-no-these/a0ac08a9-46b6-4219-8017-a0a9581ef914"
          value={getInput()}
          onInput={(e) => setInput(e.currentTarget.value)}
          class="input input-bordered w-full"
        />
        <div class="label">
          <span class="label-text-alt">{getName()}</span>
        </div>
      </label>

      <input type="submit" class="btn btn-primary" />
    </form>
  );
}
