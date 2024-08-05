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

const toTitleCase = (str: string) =>
  str[0].toUpperCase() + str.slice(1).toLowerCase();

export default function KaraCard(props: {
  kara: components["schemas"]["KaraInfoDB"];
  audioTagMap: TagMap<components["schemas"]["AudioTag"]>;
  videoTagMap: TagMap<components["schemas"]["VideoTag"]>;
}) {
  const getAudioTags = () =>
    props.kara.AudioTags?.map(
      (tag) => props.audioTagMap.get(tag.ID)?.Name,
    ).join(", ");
  const getVideoTags = () =>
    props.kara.VideoTags?.map(
      (tag) => props.videoTagMap.get(tag.ID)?.Name,
    ).join(", ");

  const getPrimary = () => {
    let primary = getAudioTags();
    const mainVideoTag = getVideoTags();
    if (primary) {
      if (props.kara.SongOrder) {
        primary += " " + props.kara.SongOrder;
      }
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
                  <a class="link-secondary">{getSourceMedia().name}</a>
                </>
              )}
            </Show>
          </p>

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
                <div class="btn btn-sm btn-ghost bg-green-700 text-base-100">
                  <FaSolidGlobe class="size-4" />
                  {getLanguage()}
                </div>
              )}
            </Show>
            <Index each={props.kara.Artists}>
              {(getArtist) => (
                <div class="btn btn-sm btn-ghost bg-amber-600 text-base-100">
                  <FaSolidMicrophoneLines class="size-4" />
                  {getArtist().Name}
                </div>
              )}
            </Index>
            <Show when={props.kara.SourceMedia}>
              {(getSourceMedia) => (
                <div class="btn btn-sm btn-ghost bg-blue-500 text-base-100">
                  <FaSolidDiagramProject class="size-4" />
                  {toTitleCase(getSourceMedia().media_type)}
                </div>
              )}
            </Show>
            <Index each={props.kara.Authors}>
              {(getAuthor) => (
                <div class="btn btn-sm btn-ghost bg-purple-600 text-base-100">
                  <FaSolidUserSecret class="size-4" />
                  {getAuthor().Name}
                </div>
              )}
            </Index>
            <div class="btn btn-sm btn-ghost bg-neutral-400 text-base-100">
              <FaSolidCalendarDays class="size-4" />
              {new Date(props.kara.CreatedAt).getFullYear()}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
