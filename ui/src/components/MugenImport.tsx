import { createEffect, createSignal, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";

export default function MugenImport(props: {
  onImport: (kara: components["schemas"]["MugenImport"]) => void;
}) {
  const [getKid, setKid] = createSignal("");
  const [getName, setName] = createSignal("");

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
      alert(resp.error);
      return;
    }

    (e.target as HTMLFormElement).reset();
    props.onImport(resp.data.import);
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full max-w-xs">
      <label>
        <div class="label">
          <span class="label-text">Karaoke ID</span>
          <span class="label-text-alt">(required)</span>
        </div>
        <input
          type="text"
          required
          value={getKid()}
          onInput={(e) => setKid(e.currentTarget.value)}
          class="input input-bordered w-full"
        />
        <div class="label">
          <span class="label-text-alt">{getName()}</span>
        </div>
      </label>

      <input type="submit" class="btn" />
    </form>
  );
}
