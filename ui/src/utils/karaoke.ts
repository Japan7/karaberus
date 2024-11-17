import Fuse from "fuse.js";
import type { components } from "./karaberus";
import type { Expression } from "fuse.js";

export type TagMap<T> = T extends { ID: infer K } ? Map<K, T> : never;

function getFuseKeys(
  audioTagMap: TagMap<components["schemas"]["AudioTag"]>,
  videoTagMap: TagMap<components["schemas"]["VideoTag"]>,
) {
  return [
    { name: "Artists.Name", weight: 2 },
    { name: "Artists.AdditionalNames.Name" },
    {
      name: "audioTagsName",
      getFn: (kara: components["schemas"]["KaraInfoDB"]) =>
        kara.AudioTags?.map((tag) => audioTagMap.get(tag.ID)?.Name ?? "") ?? [],
    },
    { name: "Authors.Name", weight: 2 },
    { name: "Comment" },
    {
      name: "creationTimeYear",
      getFn: (kara: components["schemas"]["KaraInfoDB"]) =>
        new Date(kara.KaraokeCreationTime).getFullYear().toString(),
    },
    { name: "ExtraTitles.Name" },
    { name: "Language" },
    { name: "Medias.name", weight: 2 },
    { name: "Medias.additional_name.Name" },
    { name: "SourceMedia.media_type" },
    { name: "SourceMedia.name", weight: 2 },
    { name: "Title", weight: 4 },
    { name: "Version" },
    {
      name: "videoTagsName",
      getFn: (kara: components["schemas"]["KaraInfoDB"]) =>
        kara.VideoTags?.map((tag) => videoTagMap.get(tag.ID)?.Name ?? "") ?? [],
    },
  ];
}

export function karaFuseSearch(
  karas: components["schemas"]["KaraInfoDB"][],
  audioTagMap: TagMap<components["schemas"]["AudioTag"]>,
  videoTagMap: TagMap<components["schemas"]["VideoTag"]>,
  search: string,
  tags?: { name: string; value: string }[],
) {
  const fuseKeys = getFuseKeys(audioTagMap, videoTagMap);
  const searchTerms: string[] = search.match(/("[^"]*"|[^\s]+)/g)!;

  const fuse = new Fuse(karas, {
    keys: fuseKeys,
    threshold: 0.3, // The default threshold is way too high.
    useExtendedSearch: true,
  });

  const filters: Expression[] = [];
  if (tags) {
    for (const { name, value } of tags) {
      filters.push({ [name]: `="${value}"` });
    }
  }

  // Fuse.js cannot match multiple keys on its own, so we need to iterate on
  // them. See https://github.com/krisk/Fuse/issues/235.
  if (searchTerms) {
    filters.push(
      ...searchTerms.map((term: string) => ({
        $or: fuseKeys.map((key: { name: string }) => ({
          [key.name]: term,
        })),
      })),
    );
  }

  return fuse.search({
    $and: filters,
  });
}
