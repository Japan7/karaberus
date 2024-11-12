import Fuse from "fuse.js";
import type { components } from "./karaberus";

export type TagMap<T> = T extends { ID: infer K } ? Map<K, T> : never;

export function karaFuseSearch(
  karas: components["schemas"]["KaraInfoDB"][],
  audioTagMap: TagMap<components["schemas"]["AudioTag"]>,
  videoTagMap: TagMap<components["schemas"]["VideoTag"]>,
  search: string,
) {
  const searchTerms: string[] = search.match(/("[^"]*"|[^\s]+)/g)!;

  const fuseKeys = [
    { name: "Artists.Name", weight: 2 },
    { name: "Artists.AdditionalNames.Name" },
    {
      name: "audioTagsName",
      getFn: (kara: components["schemas"]["KaraInfoDB"]) =>
        kara.AudioTags?.map((tag) => audioTagMap.get(tag.ID)?.Name ?? "") ?? [],
    },
    { name: "Authors.Name", weight: 2 },
    { name: "Comment" },
    { name: "ExtraTitles.Name" },
    { name: "Language" },
    { name: "Medias.name", weight: 2 },
    { name: "Medias.additional_name.Name" },
    { name: "SourceMedia.name", weight: 2 },
    { name: "Title", weight: 4 },
    { name: "Version" },
    {
      name: "videoTagsName",
      getFn: (kara: components["schemas"]["KaraInfoDB"]) =>
        kara.VideoTags?.map((tag) => videoTagMap.get(tag.ID)?.Name ?? "") ?? [],
    },
  ];

  const fuse = new Fuse(karas, {
    keys: fuseKeys,
    threshold: 0.3, // The default threshold is way too high.
  });

  // Fuse.js cannot match multiple keys on its own, so we need to iterate on
  // them. See https://github.com/krisk/Fuse/issues/235.
  return fuse.search({
    $and: searchTerms.map((term: string) => ({
      $or: fuseKeys.map((key: { name: string }) => ({
        [key.name]: term,
      })),
    })),
  });
}
