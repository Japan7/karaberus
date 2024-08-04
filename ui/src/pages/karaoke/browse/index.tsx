import { createResource, Index, Show } from "solid-js";
import KaraCard from "../../../components/KaraCard";
import { karaberus } from "../../../utils/karaberus-client";

export default function KaraokeBrowse() {
  const [getAllKaras] = createResource(async () => {
    const resp = await karaberus.GET("/api/kara");
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

  const getArtistMap = () =>
    new Map(getAllArtists()?.map((artist) => [artist.ID, artist]));
  const getMediaMap = () =>
    new Map(getAllMedias()?.map((media) => [media.ID, media]));
  const getAudioTagMap = () =>
    new Map(getAllAudioTags()?.map((audioTag) => [audioTag.ID, audioTag]));
  const getVideoTagMap = () =>
    new Map(getAllVideoTags()?.map((videoTag) => [videoTag.ID, videoTag]));

  return (
    <>
      <h1 class="header">Browse Karaokes</h1>

      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
        <Show when={getAllKaras()?.Karas} fallback={<p>Loading karas...</p>}>
          {(getKaras) => (
            <Index each={getKaras()}>
              {(getKara) => (
                <KaraCard
                  kara={getKara()}
                  artistMap={getArtistMap()}
                  mediaMap={getMediaMap()}
                  audioTagMap={getAudioTagMap()}
                  videoTagMap={getVideoTagMap()}
                />
              )}
            </Index>
          )}
        </Show>
      </div>
    </>
  );
}
