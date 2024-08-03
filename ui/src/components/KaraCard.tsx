import { A } from "@solidjs/router";
import {
  FaSolidCalendarDays,
  FaSolidDiagramProject,
  FaSolidGlobe,
  FaSolidMicrophoneLines,
  FaSolidUserSecret,
} from "solid-icons/fa";
import { Index, Show } from "solid-js";
import type { components } from "../utils/karaberus";

type TagMap<T> = T extends { ID: infer K } ? Map<K, T> : never;

const toTitleCase = (str: string) =>
  str[0].toUpperCase() + str.slice(1).toLowerCase();

export default function KaraCard({
  kara,
  artistMap,
  mediaMap,
  audioTagMap,
  videoTagMap,
}: {
  kara: components["schemas"]["KaraInfoDB"];
  artistMap: TagMap<components["schemas"]["Artist"]>;
  mediaMap: TagMap<components["schemas"]["MediaDB"]>;
  audioTagMap: TagMap<components["schemas"]["AudioTag"]>;
  videoTagMap: TagMap<components["schemas"]["VideoTag"]>;
}) {
  const getMainAudioTag = () => {
    const id = kara.AudioTags?.[0]?.ID;
    return (id && audioTagMap.get(id)?.Name) ?? "???";
  };
  const getMainVideoTag = () => {
    const id = kara.VideoTags?.[0]?.ID;
    return id && videoTagMap.get(id)?.Name;
  };
  const getSourceMedia = () => {
    const id = kara.SourceMedia.ID;
    return id ? mediaMap.get(id) : undefined;
  };

  return (
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">
          <A href={"./" + kara.ID} class="link-primary">
            {kara.Title}
          </A>
        </h2>

        <div class="flex flex-col gap-y-2">
          <p>
            <a class="link-secondary">{getMainAudioTag()}</a>{" "}
            <Show when={getMainVideoTag()}>
              {(getMainVideoTag) => (
                <>
                  {"("}
                  <a class="link-secondary">{getMainVideoTag()}</a>
                  {") "}
                </>
              )}
            </Show>
            from <a class="link-secondary">{getSourceMedia()?.name ?? "???"}</a>
          </p>
          <div class="flex flex-wrap gap-1">
            <Show when={kara.Language}>
              {(getLanguage) => (
                <div class="btn btn-sm btn-ghost bg-green-700">
                  <FaSolidGlobe class="size-4" />
                  {getLanguage()}
                </div>
              )}
            </Show>{" "}
            <Index each={kara.Artists}>
              {(getArtist) => (
                <div class="btn btn-sm btn-ghost bg-amber-600 text-secondary-content">
                  <FaSolidMicrophoneLines class="size-4" />
                  {artistMap.get(getArtist().ID)?.Name}
                </div>
              )}
            </Index>
            <Show when={getSourceMedia()}>
              {(getSourceMedia) => (
                <div class="btn btn-sm btn-ghost bg-blue-500 text-secondary-content">
                  <FaSolidDiagramProject class="size-4" />
                  {toTitleCase(getSourceMedia().media_type)}
                </div>
              )}
            </Show>
            <Index each={kara.Authors}>
              {(getAuthor) => (
                <div class="btn btn-sm btn-ghost bg-purple-600 text-secondary-content">
                  <FaSolidUserSecret class="size-4" />
                  {getAuthor().Name}
                </div>
              )}
            </Index>
            <div class="btn btn-sm btn-ghost bg-neutral-400 text-secondary-content">
              <FaSolidCalendarDays class="size-4" />
              {new Date(kara.CreatedAt).getFullYear()}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}