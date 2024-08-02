import { useNavigate } from "@solidjs/router";
import { createResource, createSignal, Show, type JSX } from "solid-js";
import Autocomplete from "../../components/Autocomplete";
import AutocompleteMultiple from "../../components/AutocompleteMultiple";
import type { components } from "../../utils/karaberus";
import { karaberus } from "../../utils/karaberus-client";
import routes from "../../utils/routes";

export default function KaraokeNew() {
  const navigate = useNavigate();

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
    components["schemas"]["AudioTag"][]
  >([]);
  const [getVideoTags, setVideoTags] = createSignal<
    components["schemas"]["VideoTag"][]
  >([]);
  const [getComment, setComment] = createSignal("");
  const [getVersion, setVersion] = createSignal("");
  const [getLanguage, setLanguage] = createSignal("");

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();

    const payload: components["schemas"]["KaraInfo"] = {
      title: getTitle(),
      title_aliases: getExtraTitles().trim().split("\n") || null,
      authors: getAuthors().map((author) => author.ID) || null,
      artists: getArtists().map((artist) => artist.ID) || null,
      source_media: getSourceMedia()!.ID,
      song_order: getSongOrder()!,
      medias: getMedias().map((media) => media.ID) || null,
      audio_tags: getAudioTags().map((tag) => tag.ID) || null,
      video_tags: getVideoTags().map((tag) => tag.ID) || null,
      comment: getComment(),
      version: getVersion(),
      language: getLanguage(),
    };

    const resp = await karaberus.POST("/api/kara", {
      body: payload,
    });

    if (resp.error) {
      alert(resp.error);
      return;
    }

    navigate(routes.KARAOKE_BROWSE + "/" + resp.data.kara.ID);
  };

  return (
    <>
      <h1 class="header">New Karaoke</h1>

      <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full max-w-xs">
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
            <span class="label-text">Extra titles</span>
          </div>
          <textarea
            placeholder="A Cruel Angel's Thesis&#10;Zankoku na Tenshi no Thesis"
            value={getExtraTitles()}
            oninput={(e) => setExtraTitles(e.currentTarget.value)}
            class="textarea textarea-bordered w-full"
          />
        </label>

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

        <label>
          <div class="label">
            <span class="label-text">Source media</span>
            <span class="label-text-alt">(required)</span>
          </div>
          <Show when={getAllMedias()} fallback={<p>Loading medias...</p>}>
            {(getAllMedias) => (
              <Autocomplete
                required
                placeholder="Shin Seiki Evangelion"
                items={getAllMedias()}
                getItemName={(media) => `[${media.media_type}] ${media.name}`}
                getState={getSourceMedia}
                setState={setSourceMedia}
              />
            )}
          </Show>
        </label>

        <label>
          <div class="label">
            <span class="label-text">Song order</span>
            <span class="label-text-alt">(required)</span>
          </div>
          <input
            type="number"
            min={1}
            required
            placeholder="1"
            value={getSongOrder()}
            onchange={(e) => setSongOrder(e.target.valueAsNumber)}
            class="input input-bordered w-full"
          />
        </label>

        <label>
          <div class="label">
            <span class="label-text">Medias</span>
          </div>
          <Show when={getAllMedias()} fallback={<p>Loading medias...</p>}>
            {(getAllMedias) => (
              <AutocompleteMultiple
                placeholder="Japan7"
                items={getAllMedias()}
                getItemName={(media) => `[${media.media_type}] ${media.name}`}
                getState={getMedias}
                setState={setMedias}
              />
            )}
          </Show>
        </label>

        <label>
          <div class="label">
            <span class="label-text">Audio tags</span>
          </div>
          <Show
            when={getAllAudioTags()}
            fallback={<p>Loading audio tags...</p>}
          >
            {(getAllAudioTags) => (
              <AutocompleteMultiple
                placeholder="Opening"
                items={getAllAudioTags()}
                getItemName={(tag) => tag.Name}
                getState={getAudioTags}
                setState={setAudioTags}
              />
            )}
          </Show>
        </label>

        <label>
          <div class="label">
            <span class="label-text">Video tags</span>
          </div>
          <Show
            when={getAllVideoTags()}
            fallback={<p>Loading video tags...</p>}
          >
            {(getAllVideoTags) => (
              <AutocompleteMultiple
                placeholder="Fanmade"
                items={getAllVideoTags()}
                getItemName={(tag) => tag.Name}
                getState={getVideoTags}
                setState={setVideoTags}
              />
            )}
          </Show>
        </label>

        <label>
          <div class="label">
            <span class="label-text">Comment</span>
          </div>
          <textarea
            placeholder="From https://youtu.be/dQw4w9WgXcQ"
            value={getComment()}
            oninput={(e) => setComment(e.currentTarget.value)}
            class="textarea textarea-bordered w-full"
          />
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

        <input type="submit" class="btn" />
      </form>
    </>
  );
}
