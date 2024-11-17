import { useSearchParams } from "@solidjs/router";
import { isTauri } from "@tauri-apps/api/core";
import {
  FaSolidCalendarDays,
  FaSolidDiagramProject,
  FaSolidGlobe,
  FaSolidMicrophoneLines,
  FaSolidUserSecret,
} from "solid-icons/fa";
import { Index, Show } from "solid-js";
import type { components } from "../utils/karaberus";
import type { TagMap } from "../utils/karaoke";
import MpvKaraPlayer from "./MpvKaraPlayer";

const toTitleCase = (str: string) =>
  str[0].toUpperCase() + str.slice(1).toLowerCase();

export default function KaraCard(props: {
  kara: components["schemas"]["KaraInfoDB"];
  audioTagMap: TagMap<components["schemas"]["AudioTag"]>;
  videoTagMap: TagMap<components["schemas"]["VideoTag"]>;
  enabledTagFilters: { name: string; value: string }[];
}) {
  const getAudioTags = (): string | undefined =>
    props.kara.AudioTags?.map((tag) => {
      const tagInstance = props.audioTagMap.get(tag.ID);
      if (tagInstance?.HasSongOrder && props.kara.SongOrder) {
        return `${tagInstance?.Name} ${props.kara.SongOrder}`;
      } else {
        return tagInstance?.Name;
      }
    }).join(", ");
  const getVideoTags = (): string | undefined =>
    props.kara.VideoTags?.map(
      (tag) => props.videoTagMap.get(tag.ID)?.Name,
    ).join(", ");

  const getPrimary = (): string | undefined => {
    let primary = getAudioTags();
    const mainVideoTag = getVideoTags();
    if (primary) {
      if (mainVideoTag) {
        primary += " (" + mainVideoTag + ")";
      }
    } else {
      primary = mainVideoTag;
      if (props.kara.Version) {
        primary += " (" + props.kara.Version + ")";
      }
    }
    return primary;
  };

  const [_, setSearchParams] = useSearchParams();

  function tagFilterIndexOf(name: string, value: string) {
    for (const i of props.enabledTagFilters.keys()) {
      const t = props.enabledTagFilters[i];
      if (t.name == name && t.value == value) {
        return i;
      }
    }
    return -1;
  }

  function toggleTagFilter(name: string, value: string) {
    const i = tagFilterIndexOf(name, value);
    const enabledTagFilters = props.enabledTagFilters.map((t) => t);
    if (i >= 0) {
      enabledTagFilters.splice(i, 1);
    } else {
      enabledTagFilters.push({ name: name, value: value });
    }
    setSearchParams({
      tags: enabledTagFilters
        .map((t) => `${t.name}~${encodeURIComponent(t.value)}`)
        .join(","),
      p: 0,
    });
  }

  return (
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">
          <a href={"/karaoke/browse/" + props.kara.ID} class="link-primary">
            {props.kara.Title || "[Missing Title]"}
          </a>
        </h2>

        <div class="flex flex-col gap-y-2">
          <p>
            <a class="link-secondary">{getPrimary()}</a>
            <Show when={props.kara.SourceMedia}>
              {(getSourceMedia) => (
                <>
                  {" from "}
                  <a
                    href="#"
                    class="link-secondary"
                    classList={{
                      "tag-filter-enabled":
                        tagFilterIndexOf(
                          "SourceMedia.name",
                          getSourceMedia().name,
                        ) >= 0,
                    }}
                    onclick={() =>
                      toggleTagFilter("SourceMedia.name", getSourceMedia().name)
                    }
                  >
                    {getSourceMedia().name}
                  </a>
                </>
              )}
            </Show>
          </p>

          <Show when={isTauri()}>
            <MpvKaraPlayer kara={props.kara} />
          </Show>

          <div class="flex justify-between">
            <label class="label gap-x-1">
              <input
                type="checkbox"
                disabled
                checked={props.kara.VideoUploaded}
                class="checkbox checkbox-sm checkbox-success"
              />
              <span class="label-text-alt">Video</span>
            </label>
            <label class="label gap-x-1">
              <input
                type="checkbox"
                disabled
                checked={props.kara.InstrumentalUploaded}
                class="checkbox checkbox-sm checkbox-success"
              />
              <span class="label-text-alt">Instrumental</span>
            </label>
            <label class="label gap-x-1">
              <input
                type="checkbox"
                disabled
                checked={props.kara.SubtitlesUploaded}
                class="checkbox checkbox-sm checkbox-success"
              />
              <span class="label-text-alt">Subtitles</span>
            </label>
          </div>

          <div class="flex flex-wrap gap-1">
            <Show when={props.kara.Language}>
              {(getLanguage) => (
                <div
                  class={`btn btn-sm btn-ghost text-base-100 ${
                    tagFilterIndexOf("Language", getLanguage()) >= 0
                      ? "bg-green-800 hover:bg-green-700 tag-filter-enabled"
                      : "bg-green-700 hover:bg-green-800"
                  }`}
                  onclick={() => toggleTagFilter("Language", getLanguage())}
                >
                  <FaSolidGlobe class="size-4" />
                  {getLanguage()}
                </div>
              )}
            </Show>
            <Index each={props.kara.Artists}>
              {(getArtist) => (
                <div
                  class={`btn btn-sm btn-ghost text-base-100 ${
                    tagFilterIndexOf("Artists.Name", getArtist().Name) >= 0
                      ? "bg-amber-700 hover:bg-amber-600 tag-filter-enabled"
                      : "bg-amber-600 hover:bg-amber-700"
                  }`}
                  onclick={() =>
                    toggleTagFilter("Artists.Name", getArtist().Name)
                  }
                >
                  <FaSolidMicrophoneLines class="size-4" />
                  {getArtist().Name}
                </div>
              )}
            </Index>
            <Show when={props.kara.SourceMedia}>
              {(getSourceMedia) => (
                <div
                  class={`btn btn-sm btn-ghost text-base-100 ${
                    tagFilterIndexOf(
                      "SourceMedia.media_type",
                      getSourceMedia().media_type,
                    ) >= 0
                      ? "bg-blue-600 hover:bg-blue-500 tag-filter-enabled"
                      : "bg-blue-500 hover:bg-blue-600"
                  }`}
                  onclick={() =>
                    toggleTagFilter(
                      "SourceMedia.media_type",
                      getSourceMedia().media_type,
                    )
                  }
                >
                  <FaSolidDiagramProject class="size-4" />
                  {toTitleCase(getSourceMedia().media_type)}
                </div>
              )}
            </Show>
            <Index each={props.kara.Authors}>
              {(getAuthor) => (
                <div
                  class={`btn btn-sm btn-ghost text-base-100 ${
                    tagFilterIndexOf("Authors.Name", getAuthor().Name) >= 0
                      ? "bg-purple-700 hover:bg-purple-600 tag-filter-enabled"
                      : "bg-purple-600 hover:bg-purple-700"
                  }`}
                  onclick={() =>
                    toggleTagFilter("Authors.Name", getAuthor().Name)
                  }
                >
                  <FaSolidUserSecret class="size-4" />
                  {getAuthor().Name}
                </div>
              )}
            </Index>
            <Show when={props.kara.SubtitlesUploaded}>
              <div
                class={`btn btn-sm btn-ghost text-base-100 ${
                  tagFilterIndexOf(
                    "creationTimeYear",
                    new Date(props.kara.KaraokeCreationTime)
                      .getFullYear()
                      .toString(),
                  ) >= 0
                    ? "bg-neutral-500 hover:bg-neutral-400 tag-filter-enabled"
                    : "bg-neutral-400 hover:bg-neutral-500"
                }`}
                onclick={() =>
                  toggleTagFilter(
                    "creationTimeYear",
                    new Date(props.kara.KaraokeCreationTime)
                      .getFullYear()
                      .toString(),
                  )
                }
              >
                <FaSolidCalendarDays class="size-4" />
                {new Date(props.kara.KaraokeCreationTime).getFullYear()}
              </div>
            </Show>
          </div>
        </div>
      </div>
    </div>
  );
}
