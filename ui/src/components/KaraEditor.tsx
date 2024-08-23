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
  onSubmit: (kara: components["schemas"]["KaraInfo"]) => void;
  reset?: boolean;
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

  const getAudioTag = (tagId: components["schemas"]["AudioTag"]["ID"]) =>
    (getAllAudioTags() || []).find((t) => t.ID == tagId);
  //#endregion

  //#region Signals
  const [getTitle, setTitle] = createSignal(props.kara?.Title ?? "");
  const [getExtraTitles, setExtraTitles] = createSignal(
    props.kara?.ExtraTitles?.map((v) => v.Name).join("\n") ?? "",
  );
  const [getAuthors, setAuthors] = createSignal<
    components["schemas"]["TimingAuthor"][]
  >(props.kara?.Authors ?? []);
  const [getArtists, setArtists] = createSignal<
    components["schemas"]["Artist"][]
  >(props.kara?.Artists ?? []);
  const [getSourceMedia, setSourceMedia] = createSignal<
    components["schemas"]["MediaDB"] | undefined
  >(props.kara?.SourceMedia);
  const [getSongOrder, setSongOrder] = createSignal<number | undefined>(
    props.kara?.SongOrder,
  );
  const [getMedias, setMedias] = createSignal<
    components["schemas"]["MediaDB"][]
  >(props.kara?.Medias ?? []);
  const [getAudioTags, setAudioTags] = createSignal<
    components["schemas"]["AudioTagDB"][]
  >(props.kara?.AudioTags ?? []);
  const [getVideoTags, setVideoTags] = createSignal<
    components["schemas"]["VideoTagDB"][]
  >(props.kara?.VideoTags ?? []);
  const [getComment, setComment] = createSignal(props.kara?.Comment ?? "");
  const [getVersion, setVersion] = createSignal(props.kara?.Version ?? "");
  const [getLanguage, setLanguage] = createSignal(props.kara?.Language ?? "");
  //#endregion

  //#region Handlers
  let modalRef!: HTMLDialogElement;
  const [getModalForm, setModalForm] = createSignal<JSX.Element>();

  const [getToast, setToast] = createSignal<string>();

  const showToast = (msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(), 3000);
  };

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = (e) => {
    e.preventDefault();
    props.onSubmit({
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
    });
    if (props.reset) {
      (e.target as HTMLFormElement).reset();
    }
  };

  const postAuthor = async (author: components["schemas"]["AuthorInfo"]) => {
    const resp = await karaberus.POST("/api/tags/author", { body: author });
    if (resp.error) {
      alert(resp.error.title);
      return;
    }
    showToast("Author added!");
    refetchAuthors();
    modalRef.close();
  };

  const openAddAuthorModal: JSX.EventHandler<HTMLElement, MouseEvent> = (e) => {
    e.preventDefault();
    setModalForm(<AuthorEditor onSubmit={postAuthor} />);
    modalRef.showModal();
  };

  const postArtist = async (artist: components["schemas"]["ArtistInfo"]) => {
    const resp = await karaberus.POST("/api/tags/artist", { body: artist });
    if (resp.error) {
      alert(resp.error.title);
      return;
    }
    showToast("Artist added!");
    refetchArtists();
    modalRef.close();
  };

  const openAddArtistModal: JSX.EventHandler<HTMLElement, MouseEvent> = (e) => {
    e.preventDefault();
    setModalForm(<ArtistEditor onSubmit={postArtist} />);
    modalRef.showModal();
  };

  const postMedia = async (media: components["schemas"]["MediaInfo"]) => {
    const resp = await karaberus.POST("/api/tags/media", { body: media });
    if (resp.error) {
      alert(resp.error.title);
      return;
    }
    showToast("Media added!");
    refetchMedia();
    modalRef.close();
  };

  const openAddMediaModal: JSX.EventHandler<HTMLElement, MouseEvent> = (e) => {
    e.preventDefault();
    setModalForm(<MediaEditor onSubmit={postMedia} />);
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
        placeholder="Zankoku na Tenshi no These"
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
        placeholder={"Zankoku na Tenshi no Thesis\nA Cruel Angel's Thesis"}
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
            Can't find it?
          </a>
        </span>
      </div>
      <Show
        when={getAllAuthors()}
        fallback={<span class="loading loading-spinner loading-lg" />}
      >
        {(getAllAuthors) => (
          <AutocompleteMultiple
            items={getAllAuthors()}
            getItemName={(author) => author.Name}
            getState={getAuthors}
            setState={setAuthors}
            placeholder="bebou69"
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
            Can't find it?
          </a>
        </span>
      </div>
      <Show
        when={getAllArtists()}
        fallback={<span class="loading loading-spinner loading-lg" />}
      >
        {(getAllArtists) => (
          <AutocompleteMultiple
            items={getAllArtists()}
            getItemName={(artist) => artist.Name}
            getState={getArtists}
            setState={setArtists}
            placeholder="Yoko Takahashi"
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
            Can't find it?
          </a>
        </span>
      </div>
      <Show
        when={getAllMedias()}
        fallback={<span class="loading loading-spinner loading-lg" />}
      >
        {(getAllMedias) => (
          <Autocomplete
            items={getAllMedias()}
            getItemName={(media) => `[${media.media_type}] ${media.name}`}
            getState={getSourceMedia}
            setState={setSourceMedia}
            placeholder="Shin Seiki Evangelion"
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
        placeholder="1"
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
            Can't find it?
          </a>
        </span>
      </div>
      <Show
        when={getAllMedias()}
        fallback={<span class="loading loading-spinner loading-lg" />}
      >
        {(getAllMedias) => (
          <AutocompleteMultiple
            items={getAllMedias()}
            getItemName={(media) => `[${media.media_type}] ${media.name}`}
            getState={getMedias}
            setState={setMedias}
            placeholder="Japan7"
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
        <Index
          each={getAllAudioTags()}
          fallback={<span class="loading loading-spinner loading-lg" />}
        >
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
        <Index
          each={getAllVideoTags()}
          fallback={<span class="loading loading-spinner loading-lg" />}
        >
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
        placeholder="something something"
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
        placeholder="gael42"
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
        placeholder="jpn"
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
        <input type="submit" class="btn btn-primary" />
      </form>

      <dialog ref={modalRef} class="modal modal-bottom sm:modal-middle">
        <div class="modal-box flex justify-center">{getModalForm()}</div>
        <form method="dialog" class="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>

      <Show when={getToast()}>
        {(getToast) => (
          <div class="toast">
            <div class="alert alert-success">
              <span>{getToast()}</span>
            </div>
          </div>
        )}
      </Show>
    </>
  );
}
