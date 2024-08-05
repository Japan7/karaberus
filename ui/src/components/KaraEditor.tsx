import { createResource, createSignal, Index, Show, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";
import ArtistEditor from "./ArtistEditor";
import AuthorEditor from "./AuthorEditor";
import Autocomplete from "./Autocomplete";
import AutocompleteMultiple from "./AutocompleteMultiple";
import MediaEditor from "./MediaEditor";

export default function KaraEditor(props: {
  kara?: components["schemas"]["KaraInfoDB"];
  onSubmit: (info: components["schemas"]["KaraInfo"]) => void;
}) {
  //#region Resources
  const [getAllAuthors, { refetch: refetchAuthors }] = createResource(
    async () => {
      const resp = await karaberus.GET("/api/tags/author");
      return resp.data;
    },
  );
  const [getAllArtists, { refetch: refetchArtists }] = createResource(
    async () => {
      const resp = await karaberus.GET("/api/tags/artist");
      return resp.data;
    },
  );
  const [getAllMedias, { refetch: refetchMedia }] = createResource(async () => {
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

  const getAudioTag = (tagId: string) =>
    (getAllAudioTags() || []).find((t) => t.ID == tagId);
  //#endregion

  //#region Signals
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

  if (props.kara) {
    setTitle(props.kara.Title);
    if (props.kara.ExtraTitles) {
      setExtraTitles(props.kara.ExtraTitles.map((v) => v.Name).join("\n"));
    }
    if (props.kara.Authors) {
      setAuthors(props.kara.Authors);
    }
    if (props.kara.Artists) {
      setArtists(props.kara.Artists);
    }
    if (props.kara.SourceMedia) {
      setSourceMedia(props.kara.SourceMedia);
    }
    setSongOrder(props.kara.SongOrder);
    if (props.kara.Medias) {
      setMedias(props.kara.Medias);
    }
    if (props.kara.AudioTags) {
      setAudioTags(props.kara.AudioTags);
    }
    if (props.kara.VideoTags) {
      setVideoTags(props.kara.VideoTags);
    }
    setComment(props.kara.Comment);
    setVersion(props.kara.Version);
    setLanguage(props.kara.Language);
  }
  //#endregion

  //#region Handlers
  let modalRef!: HTMLDialogElement;

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

    props.onSubmit(payload);
  };

  const [getModalForm, setModalForm] = createSignal<JSX.Element>();

  const openAddAuthorModal: JSX.EventHandler<HTMLElement, MouseEvent> = (e) => {
    e.preventDefault();
    setModalForm(
      <AuthorEditor
        onAdd={() => {
          refetchAuthors();
          modalRef.close();
        }}
      />,
    );
    modalRef.showModal();
  };

  const openAddArtistModal: JSX.EventHandler<HTMLElement, MouseEvent> = (e) => {
    e.preventDefault();
    setModalForm(
      <ArtistEditor
        onAdd={() => {
          refetchArtists();
          modalRef.close();
        }}
      />,
    );
    modalRef.showModal();
  };

  const openAddMediaModal: JSX.EventHandler<HTMLElement, MouseEvent> = (e) => {
    e.preventDefault();
    setModalForm(
      <MediaEditor
        onAdd={() => {
          refetchMedia();
          modalRef.close();
        }}
      />,
    );
    modalRef.showModal();
  };
  //#endregion

  //#region Render
  const titleInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Title</span>
        <span class="label-text-alt">(required)</span>
      </div>
      <input
        type="text"
        required
        value={getTitle()}
        oninput={(e) => setTitle(e.currentTarget.value)}
        class="input input-bordered w-full"
      />
    </label>
  );
  const extraTitlesInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Aliases</span>
        <span class="label-text-alt">1 per line</span>
      </div>
      <textarea
        value={getExtraTitles()}
        oninput={(e) => setExtraTitles(e.currentTarget.value)}
        class="textarea textarea-bordered w-full"
      />
    </label>
  );
  const authorsInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Authors</span>
        <span class="label-text-alt">
          <a class="link" onclick={openAddAuthorModal}>
            Add new
          </a>
        </span>
      </div>
      <Show when={getAllAuthors()} fallback={<p>Loading Authors...</p>}>
        {(getAllAuthors) => (
          <AutocompleteMultiple
            items={getAllAuthors()}
            getItemName={(author) => author.Name}
            getState={getAuthors}
            setState={setAuthors}
          />
        )}
      </Show>
    </label>
  );
  const artistsInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Artists</span>
        <span class="label-text-alt">
          <a class="link" onclick={openAddArtistModal}>
            Add new
          </a>
        </span>
      </div>
      <Show when={getAllArtists()} fallback={<p>Loading artists...</p>}>
        {(getAllArtists) => (
          <AutocompleteMultiple
            items={getAllArtists()}
            getItemName={(artist) => artist.Name}
            getState={getArtists}
            setState={setArtists}
          />
        )}
      </Show>
    </label>
  );
  const sourceMediaInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Source media</span>
        <span class="label-text-alt">
          <a class="link" onclick={openAddMediaModal}>
            Add new
          </a>
        </span>
      </div>
      <Show when={getAllMedias()} fallback={<p>Loading medias...</p>}>
        {(getAllMedias) => (
          <Autocomplete
            items={getAllMedias()}
            getItemName={(media) => `[${media.media_type}] ${media.name}`}
            getState={getSourceMedia}
            setState={setSourceMedia}
          />
        )}
      </Show>
    </label>
  );
  const songOrderInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Song order</span>
      </div>
      <input
        type="number"
        min={0}
        value={getSongOrder()}
        onchange={(e) => setSongOrder(e.target.valueAsNumber)}
        class="input input-bordered w-full"
      />
    </label>
  );
  const mediasInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Other medias</span>
        <span class="label-text-alt">
          <a class="link" onclick={openAddMediaModal}>
            Add new
          </a>
        </span>
      </div>
      <Show when={getAllMedias()} fallback={<p>Loading medias...</p>}>
        {(getAllMedias) => (
          <AutocompleteMultiple
            items={getAllMedias()}
            getItemName={(media) => `[${media.media_type}] ${media.name}`}
            getState={getMedias}
            setState={setMedias}
          />
        )}
      </Show>
    </label>
  );
  const audioTagsInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Song types</span>
      </div>
      <div class="grid md:grid-cols-2">
        <Index each={getAllAudioTags()} fallback={<p>Loading audio tags...</p>}>
          {(getAudioTag) => (
            <label class="label cursor-pointer justify-start gap-x-2">
              <input
                type="checkbox"
                checked={getAudioTags().some(
                  (tag) => tag.ID == getAudioTag().ID,
                )}
                onchange={(e) =>
                  setAudioTags((prev) =>
                    e.currentTarget.checked
                      ? [...prev, getAudioTag()]
                      : prev.filter((tag) => tag.ID != getAudioTag().ID),
                  )
                }
                class="checkbox"
              />
              <span class="label-text">{getAudioTag().Name}</span>
            </label>
          )}
        </Index>
      </div>
    </label>
  );
  const videoTagsInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Tags</span>
      </div>
      <div class="grid md:grid-cols-2">
        <Index each={getAllVideoTags()} fallback={<p>Loading video tags...</p>}>
          {(getVideoTag) => (
            <label class="label cursor-pointer justify-start gap-x-2">
              <input
                type="checkbox"
                checked={getVideoTags().some(
                  (tag) => tag.ID == getVideoTag().ID,
                )}
                onchange={(e) =>
                  setVideoTags((prev) =>
                    e.currentTarget.checked
                      ? [...prev, getVideoTag()]
                      : prev.filter((tag) => tag.ID != getVideoTag().ID),
                  )
                }
                class="checkbox"
              />
              <span class="label-text">{getVideoTag().Name}</span>
            </label>
          )}
        </Index>
      </div>
    </label>
  );
  const commentInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Comment</span>
      </div>
      <textarea
        value={getComment()}
        oninput={(e) => setComment(e.currentTarget.value)}
        class="textarea textarea-bordered w-full"
      />
    </label>
  );
  const versionInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Version</span>
      </div>
      <input
        type="text"
        value={getVersion()}
        oninput={(e) => setVersion(e.currentTarget.value)}
        class="input input-bordered w-full"
      />
    </label>
  );
  const languageInput = () => (
    <label>
      <div class="label">
        <span class="label-text">Language</span>
      </div>
      <input
        type="text"
        value={getLanguage()}
        oninput={(e) => setLanguage(e.currentTarget.value)}
        class="input input-bordered w-full"
      />
    </label>
  );
  //#endregion

  return (
    <>
      <form onsubmit={onsubmit} class="flex flex-col gap-y-4">
        <div class="grid md:grid-cols-2 gap-4">
          <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
              <h2 class="card-title">Titles</h2>
              {titleInput()}
              {extraTitlesInput()}
            </div>
          </div>
          <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
              <h2 class="card-title">Audio</h2>
              {artistsInput()}
              {audioTagsInput()}
              <Show
                when={getAudioTags().some(
                  (tag) => getAudioTag(tag.ID)?.HasSongOrder,
                )}
              >
                {sourceMediaInput()}
                {songOrderInput()}
              </Show>
              {languageInput()}
            </div>
          </div>
          <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
              <h2 class="card-title">Video</h2>
              {videoTagsInput()}
              {mediasInput()}
            </div>
          </div>
          <div class="card bg-base-100 shadow-xl md:row-start-1 md:col-start-2">
            <div class="card-body">
              <h2 class="card-title">Additional informations</h2>
              {authorsInput()}
              {commentInput()}
              {versionInput()}
            </div>
          </div>
        </div>
        <input type="submit" class="btn" />
      </form>

      <dialog ref={modalRef} class="modal modal-bottom sm:modal-middle">
        <div class="modal-box flex justify-center">{getModalForm()}</div>
        <form method="dialog" class="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>
    </>
  );
}
