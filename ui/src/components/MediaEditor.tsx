import { createResource, createSignal, Index, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";

export default function MediaEditor(props: {
  onAdd: (media: components["schemas"]["MediaDB"]) => void;
}) {
  const [getAllMediaTypes] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/media/types");
    return resp.data;
  });

  const [getMediaType, setMediaType] = createSignal("ANIME");
  const [getName, setName] = createSignal("");
  const [getAdditionalNames, setAdditionalNames] = createSignal("");

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();
    const resp = await karaberus.POST("/api/tags/media", {
      body: {
        media_type: getMediaType(),
        name: getName(),
        additional_names: getAdditionalNames().trim().split("\n"),
      },
    });
    if (resp.error) {
      alert(resp.error);
      return;
    }
    (e.target as HTMLFormElement).reset();
    props.onAdd(resp.data.media);
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full max-w-xs">
      <label>
        <div class="label">
          <span class="label-text">Media type</span>
          <span class="label-text-alt">(required)</span>
        </div>
        <select
          required
          value={getMediaType()}
          onchange={(e) => setMediaType(e.currentTarget.value)}
          class="select select-bordered w-full"
        >
          <Index each={getAllMediaTypes()}>
            {(getMediaType) => (
              <option value={getMediaType().ID}>{getMediaType().Name}</option>
            )}
          </Index>
        </select>
      </label>

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

      <label>
        <div class="label">
          <span class="label-text">Additional names</span>
          <span class="label-text-alt">1 per line</span>
        </div>
        <textarea
          value={getAdditionalNames()}
          onInput={(e) => setAdditionalNames(e.currentTarget.value)}
          class="textarea textarea-bordered w-full"
        />
      </label>

      <input type="submit" class="btn" />
    </form>
  );
}
