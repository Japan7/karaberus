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
  Show,
} from "solid-js";
import KaraCard from "../../../components/KaraCard";
import { karaberus } from "../../../utils/karaberus-client";
import { karaFuseSearch } from "../../../utils/karaoke";

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

  const [getSearch, setSearch] = createSignal("");
  const [getPage, setPage] = createSignal(0);

  createEffect(() => {
    getSearch();
    setPage(0);
  });

  const getAudioTagMap = () =>
    new Map(getAllAudioTags()?.map((audioTag) => [audioTag.ID, audioTag]));
  const getVideoTagMap = () =>
    new Map(getAllVideoTags()?.map((videoTag) => [videoTag.ID, videoTag]));

  const getFilteredKaras = () => {
    const karas = getAllKaras();
    if (!karas) {
      return [];
    }

    const search = getSearch();
    if (!search) {
      return karas;
    }

    const results = karaFuseSearch(
      karas,
      getAudioTagMap(),
      getVideoTagMap(),
      search,
    );

    return results.map((result) => result.item);
  };

  const pageSize = 9;

  const getPageKaras = () => {
    const karas = getFilteredKaras();
    const page = getPage();
    const start = page * pageSize;
    const end = start + pageSize;
    return karas.slice(start, end);
  };

  return (
    <>
      <h1 class="header">Browse Karaokes</h1>

      <label class="input input-bordered flex items-center gap-2">
        <input
          type="text"
          placeholder="Search"
          value={getSearch()}
          oninput={(e) => setSearch(e.currentTarget.value)}
          class="grow"
        />
        <HiOutlineMagnifyingGlass class="size-4 opacity-70" />
      </label>

      <div class="join mx-auto">
        <button
          disabled={getPage() === 0}
          onclick={() => setPage((page) => page - 1)}
          class="join-item btn"
        >
          <HiSolidChevronDoubleLeft class="size-5" />
        </button>
        <button onclick={() => setPage(0)} class="join-item btn">
          Page {getPage() + 1}
        </button>
        <button
          disabled={
            getPage() * pageSize + pageSize >= getFilteredKaras().length
          }
          onclick={() => setPage((page) => page + 1)}
          class="join-item btn"
        >
          <HiSolidChevronDoubleRight class="size-5" />
        </button>
      </div>

      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
        <Show when={getAllKaras()} fallback={<p>Loading karas...</p>}>
          <Index each={getPageKaras()} fallback={<p>No results</p>}>
            {(getKara) => (
              <KaraCard
                kara={getKara()}
                audioTagMap={getAudioTagMap()}
                videoTagMap={getVideoTagMap()}
              />
            )}
          </Index>
        </Show>
      </div>
    </>
  );
}
