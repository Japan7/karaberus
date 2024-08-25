import { createResource, createSignal, Index, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";

export default function MediaEditor(props: {
  media?: components["schemas"]["MediaDB"];
  onSubmit: (media: components["schemas"]["MediaInfo"]) => void;
  reset?: boolean;
}) {
  const [getAllMediaTypes] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/media/types");
    return resp.data;
  });

  const [getMediaType, setMediaType] = createSignal(
    props.media?.media_type ?? "ANIME",
  );
  const [getName, setName] = createSignal(props.media?.name ?? "");
  const [getAdditionalNames, setAdditionalNames] = createSignal(
    props.media?.additional_name?.map((n) => n.Name).join("\n") ?? "",
  );

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = (e) => {
    e.preventDefault();
    const additionalNamesStr = getAdditionalNames().trim();
    const additionalNames = additionalNamesStr
      ? additionalNamesStr.split("\n")
      : null;
    props.onSubmit({
      media_type: getMediaType(),
      name: getName(),
      additional_names: additionalNames,
    });
    if (props.reset) {
      (e.target as HTMLFormElement).reset();
    }
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full">
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
            {(getOptionMediaType) => (
              <option
                value={getOptionMediaType().ID}
                selected={getMediaType() === getOptionMediaType().ID}
              >
                {getOptionMediaType().Name}
              </option>
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
          placeholder="Shin Seiki Evangelion"
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
          placeholder={"Neon Genesis Evangelion\n新世紀エヴァンゲリオン"}
          value={getAdditionalNames()}
          onInput={(e) => setAdditionalNames(e.currentTarget.value)}
          class="textarea textarea-bordered w-full"
        />
      </label>

      <input type="submit" class="btn btn-primary" />
    </form>
  );
}
