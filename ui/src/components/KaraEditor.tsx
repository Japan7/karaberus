import { createResource, createSignal, Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";
import Autocomplete from "./Autocomplete";
import AutocompleteMultiple from "./AutocompleteMultiple";

function getAudioTag(
  allAudioTags: components["schemas"]["AudioTag"][],
  tagId: string,
) {
  return allAudioTags.find((t) => t.ID == tagId);
}

function getVideoTag(
  allVideoTags: components["schemas"]["VideoTag"][],
  tagId: string,
) {
  return allVideoTags.find((t) => t.ID == tagId);
}

export default function KaraEditor({
  kara,
  onSubmit,
}: {
  kara?: components["schemas"]["KaraInfoDB"];
  onSubmit: (info: components["schemas"]["KaraInfo"]) => void;
}) {
  const [getAllAuthors] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/author");
    return resp.data;
  });
  const [getAllArtists] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/artist");
    return resp.data;
  });
  const [getAllMedias] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/media");
    return resp.data;
  });
  const [getAllAudioTags] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/audio");
    return resp.data;
  });
  const [getAllVideoTags] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/video");
    return resp.data;
  });

  const [getTitle, setTitle] = createSignal("");
  const [getExtraTitles, setExtraTitles] = createSignal("");
  const [getAuthors, setAuthors] = createSignal<
    components["schemas"]["TimingAuthor"][]
  >([]);
  const [getArtists, setArtists] = createSignal<
    components["schemas"]["Artist"][]
  >([]);
  const [getSourceMedia, setSourceMedia] =
    createSignal<components["schemas"]["MediaDB"]>();
  const [getSongOrder, setSongOrder] = createSignal<number>();
  const [getMedias, setMedias] = createSignal<
    components["schemas"]["MediaDB"][]
  >([]);
  const [getAudioTags, setAudioTags] = createSignal<
    components["schemas"]["AudioTagDB"][]
  >([]);
  const [getVideoTags, setVideoTags] = createSignal<
    components["schemas"]["VideoTagDB"][]
  >([]);
  const [getComment, setComment] = createSignal("");
  const [getVersion, setVersion] = createSignal("");
  const [getLanguage, setLanguage] = createSignal("");

  if (kara) {
    console.log(kara);
    setTitle(kara.Title);
    if (kara.ExtraTitles)
      setExtraTitles(kara.ExtraTitles.map((v) => v.Name).join("\n"));

    if (kara.Authors) setAuthors(kara.Authors);

    if (kara.SourceMedia) setSourceMedia(kara.SourceMedia);
    if (kara.Medias) setMedias(kara.Medias);

    //tags
    if (kara.AudioTags) setAudioTags(kara.AudioTags);
    if (kara.VideoTags) setVideoTags(kara.VideoTags);

    setComment(kara.Comment);
    setVersion(kara.Version);
    setLanguage(kara.Language);

    setSongOrder(kara.SongOrder);
  }

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();

    const payload: components["schemas"]["KaraInfo"] = {
      title: getTitle(),
      title_aliases: getExtraTitles().trim().split("\n") || null,
      authors: getAuthors().map((author) => author.ID) || null,
      artists: getArtists().map((artist) => artist.ID) || null,
      source_media: getSourceMedia()?.ID || 0,
      song_order: getSongOrder() || 0,
      medias: getMedias().map((media) => media.ID) || null,
      audio_tags: getAudioTags().map((tag) => tag.ID) || null,
      video_tags: getVideoTags().map((tag) => tag.ID) || null,
      comment: getComment(),
      version: getVersion(),
      language: getLanguage(),
    };

    onSubmit(payload);
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-4">
      <div class="grid md:grid-cols-2 gap-4">
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h2 class="card-title">Title(s)</h2>
            <label>
              <div class="label">
                <span class="label-text">Title</span>
                <span class="label-text-alt">(required)</span>
              </div>
              <input
                type="text"
                required
                placeholder="Zankoku na Tenshi no These"
                value={getTitle()}
                oninput={(e) => setTitle(e.currentTarget.value)}
                class="input input-bordered w-full"
              />
            </label>
            <label>
              <div class="label">
                <span class="label-text">Aliases</span>
              </div>
              <textarea
                placeholder="A Cruel Angel's Thesis&#10;Zankoku na Tenshi no Thesis"
                value={getExtraTitles()}
                oninput={(e) => setExtraTitles(e.currentTarget.value)}
                class="textarea textarea-bordered w-full"
              />
            </label>
          </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h2 class="card-title">Identity</h2>
            <label>
              <div class="label">
                <span class="label-text">Language</span>
              </div>
              <input
                type="text"
                placeholder="FR"
                value={getLanguage()}
                oninput={(e) => setLanguage(e.currentTarget.value)}
                class="input input-bordered w-full"
              />
            </label>
            <label>
              <div class="label">
                <span class="label-text">Source media</span>
              </div>
              <Show when={getAllMedias()} fallback={<p>Loading medias...</p>}>
                {(getAllMedias) => (
                  <Autocomplete
                    placeholder="Shin Seiki Evangelion"
                    items={getAllMedias()}
                    getItemName={(media) =>
                      `[${media.media_type}] ${media.name}`
                    }
                    getState={getSourceMedia}
                    setState={setSourceMedia}
                  />
                )}
              </Show>
            </label>
            <label>
              <div class="label">
                <span class="label-text">Other medias</span>
              </div>
              <Show when={getAllMedias()} fallback={<p>Loading medias...</p>}>
                {(getAllMedias) => (
                  <AutocompleteMultiple
                    placeholder="Japan7"
                    items={getAllMedias()}
                    getItemName={(media) =>
                      `[${media.media_type}] ${media.name}`
                    }
                    getState={getMedias}
                    setState={setMedias}
                  />
                )}
              </Show>
            </label>
            <label>
              <div class="label">
                <span class="label-text">Song types</span>
              </div>
              <Show
                when={getAllAudioTags()}
                fallback={<p>Loading audio tags...</p>}
              >
                {(getAllAudioTags) => (
                  <AutocompleteMultiple
                    placeholder="Opening"
                    items={getAllAudioTags()}
                    getItemName={(tag) =>
                      getAudioTag(getAllAudioTags(), tag.ID)?.Name || ""
                    }
                    getState={getAudioTags}
                    setState={setAudioTags}
                  />
                )}
              </Show>
            </label>
            <Show
              when={getAudioTags().some(
                (tag) =>
                  getAudioTag(getAllAudioTags() || [], tag.ID)?.HasSongOrder,
              )}
            >
              <label>
                <div class="label">
                  <span class="label-text">Song order</span>
                </div>
                <input
                  type="number"
                  min={0}
                  placeholder="1"
                  value={getSongOrder()}
                  onchange={(e) => setSongOrder(e.target.valueAsNumber)}
                  class="input input-bordered w-full"
                />
              </label>
            </Show>
            <label>
              <div class="label">
                <span class="label-text">Artists</span>
              </div>
              <Show when={getAllArtists()} fallback={<p>Loading artists...</p>}>
                {(getAllArtists) => (
                  <AutocompleteMultiple
                    placeholder="Yoko Takahashi"
                    items={getAllArtists()}
                    getItemName={(artist) => artist.Name}
                    getState={getArtists}
                    setState={setArtists}
                  />
                )}
              </Show>
            </label>
          </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h2 class="card-title">Categorization</h2>
            <label>
              <div class="label">
                <span class="label-text">Tags</span>
              </div>
              <Show
                when={getAllVideoTags()}
                fallback={<p>Loading video tags...</p>}
              >
                {(getAllVideoTags) => (
                  <AutocompleteMultiple
                    placeholder="Fanmade"
                    items={getAllVideoTags()}
                    getItemName={(tag) =>
                      getVideoTag(getAllVideoTags(), tag.ID)?.Name || ""
                    }
                    getState={getVideoTags}
                    setState={setVideoTags}
                  />
                )}
              </Show>
            </label>
            <label>
              <div class="label">
                <span class="label-text">Version</span>
              </div>
              <input
                type="text"
                placeholder="iykyk"
                value={getVersion()}
                oninput={(e) => setVersion(e.currentTarget.value)}
                class="input input-bordered w-full"
              />
            </label>
          </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h2 class="card-title">Metadata</h2>
            <label>
              <div class="label">
                <span class="label-text">Authors</span>
              </div>
              <Show when={getAllAuthors()} fallback={<p>Loading Authors...</p>}>
                {(getAllAuthors) => (
                  <AutocompleteMultiple
                    placeholder="RhiobeT"
                    items={getAllAuthors()}
                    getItemName={(author) => author.Name}
                    getState={getAuthors}
                    setState={setAuthors}
                  />
                )}
              </Show>
            </label>
            <label>
              <div class="label">
                <span class="label-text">Comments</span>
              </div>
              <textarea
                placeholder="From https://youtu.be/dQw4w9WgXcQ"
                value={getComment()}
                oninput={(e) => setComment(e.currentTarget.value)}
                class="textarea textarea-bordered w-full"
              />
            </label>
          </div>
        </div>
      </div>

      <input type="submit" class="btn" />
    </form>
  );
}
