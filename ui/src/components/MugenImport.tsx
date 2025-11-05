import { createEffect, createSignal, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";

interface MugenKaraTitles {
  [key: string]: string;
}

interface MugenKara {
  titles: MugenKaraTitles;
  titles_default_language: string;
}

async function getImportTitle(kid: string): Promise<string> {
  if (kid.length !== 36) {
    return "";
  }
  const resp = await fetch(`https://kara.moe/api/karas/${kid}`);
  const json: MugenKara = await resp.json();

  // shouldn’t be undefined but who knows
  return (
    json.titles[json.titles_default_language] ?? "couldn’t get import title"
  );
}

export default function MugenImport(props: {
  onImport: (kara: components["schemas"]["MugenImport"]) => void;
}) {
  const [getInput, setInput] = createSignal("");
  const [getName, setName] = createSignal("");

  const getKid = () => {
    const input = getInput().trim();
    return input.split("/").pop() ?? input;
  };

  createEffect(async () => {
    setName(await getImportTitle(getKid()));
  });

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();

    const resp = await karaberus.POST("/api/mugen", {
      body: {
        mugen_kid: getKid(),
      },
    });

    if (resp.response.status === 204) {
      alert("karaoke already imported");
      return;
    }

    if (resp.error) {
      alert(resp.error.detail);
      return;
    }

    (e.target as HTMLFormElement).reset();
    props.onImport(resp.data.import);
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full">
      <label>
        <div class="label">
          <span class="">Karaoke ID/URL</span>
          <span class="text-sm opacity-70">(required)</span>
        </div>
        <input
          type="text"
          required
          placeholder="https://kara.moe/kara/zankoku-na-tenshi-no-these/a0ac08a9-46b6-4219-8017-a0a9581ef914"
          value={getInput()}
          onInput={(e) => setInput(e.currentTarget.value)}
          class="input input w-full"
        />
        <div class="label">
          <span class="text-sm opacity-70">{getName()}</span>
        </div>
      </label>

      <input type="submit" class="btn btn-primary" />
    </form>
  );
}
