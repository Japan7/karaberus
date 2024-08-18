import { debounce } from "@solid-primitives/scheduled";
import {
  HiOutlineMagnifyingGlass,
  HiSolidChevronDoubleLeft,
  HiSolidChevronDoubleRight,
} from "solid-icons/hi";
import {
  createEffect,
  createResource,
  createSignal,
  Index,
  on,
  Show,
} from "solid-js";
import KaraCard from "../../../components/KaraCard";
import type { components } from "../../../utils/karaberus";
import { karaberus } from "../../../utils/karaberus-client";
import { karaFuseSearch } from "../../../utils/karaoke";

const PAGE_SIZE = 9;

export default function KaraokeBrowse() {
  const [getAllKaras] = createResource(async () => {
    const resp = await karaberus.GET("/api/kara");
    return resp.data?.Karas?.sort((a, b) => b.ID - a.ID);
  });
  const [getAllAudioTags] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/audio");
    return resp.data;
  });
  const [getAllVideoTags] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/video");
    return resp.data;
  });

  const [getPage, setPage] = createSignal(0);

  const [getQuery, setQuery] = createSignal("");
  const [getSearchResults, setSearchResults] = createSignal<
    components["schemas"]["KaraInfoDB"][]
  >([]);

  createEffect(
    on(getAllKaras, () => {
      debouncedSearch(getQuery());
    }),
  );

  createEffect(
    on(getSearchResults, () => {
      setPage(0);
    }),
  );

  const getAudioTagMap = () =>
    new Map(getAllAudioTags()?.map((audioTag) => [audioTag.ID, audioTag]));
  const getVideoTagMap = () =>
    new Map(getAllVideoTags()?.map((videoTag) => [videoTag.ID, videoTag]));

  const search = (query: string) => {
    const karas = getAllKaras();
    if (!karas) {
      return [];
    }

    if (!query) {
      return karas;
    }

    const results = karaFuseSearch(
      karas,
      getAudioTagMap(),
      getVideoTagMap(),
      query,
    );

    return results.map((result) => result.item);
  };
  const debouncedSearch = debounce((query: string) => {
    setSearchResults(search(query));
  }, 250);

  createEffect(() => {
    debouncedSearch(getQuery());
  });

  const getTotal = () => getSearchResults().length;

  const getStart = () => getPage() * PAGE_SIZE;
  const getEnd = () => Math.min(getStart() + PAGE_SIZE, getTotal());

  const getPageKaras = () => getSearchResults().slice(getStart(), getEnd());

  const nav = () => (
    <div class="join mx-auto">
      <button
        disabled={getPage() === 0}
        onclick={() => setPage((page) => page - 1)}
        class="join-item btn"
      >
        <HiSolidChevronDoubleLeft class="size-5" />
      </button>
      <select
        onchange={(e) => setPage(parseInt(e.currentTarget.value))}
        class="join-item select bg-base-200"
      >
        <Index
          each={[
            ...Array(Math.ceil(getSearchResults().length / PAGE_SIZE)).keys(),
          ]}
        >
          {(getIndex) => (
            <option value={getIndex()} selected={getIndex() === getPage()}>
              Page {getIndex() + 1}
            </option>
          )}
        </Index>
      </select>
      <button
        disabled={
          getPage() * PAGE_SIZE + PAGE_SIZE >= getSearchResults().length
        }
        onclick={() => setPage((page) => page + 1)}
        class="join-item btn"
      >
        <HiSolidChevronDoubleRight class="size-5" />
      </button>
    </div>
  );

  return (
    <>
      <h1 class="header">Browse Karaokes</h1>

      <Show
        when={getAllKaras()}
        fallback={<span class="loading loading-spinner loading-lg" />}
      >
        <label class="input input-bordered flex items-center gap-2">
          <input
            type="text"
            placeholder="Search"
            value={getQuery()}
            oninput={(e) => setQuery(e.currentTarget.value)}
            class="grow"
          />
          <HiOutlineMagnifyingGlass class="size-4 opacity-70" />
        </label>

        <Show when={getSearchResults().length} fallback={<p>No results</p>}>
          {nav()}

          <p>
            Showing <b>{getStart() + 1}</b> to <b>{getEnd()}</b> of{" "}
            <b>{getTotal()}</b> results
          </p>

          <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            <Index each={getPageKaras()}>
              {(getKara) => (
                <KaraCard
                  kara={getKara()}
                  audioTagMap={getAudioTagMap()}
                  videoTagMap={getVideoTagMap()}
                />
              )}
            </Index>
          </div>

          {nav()}
        </Show>
      </Show>
    </>
  );
}
