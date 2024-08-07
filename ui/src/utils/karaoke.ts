import Fuse from "fuse.js";
import type { components } from "./karaberus";

export type TagMap<T> = T extends { ID: infer K } ? Map<K, T> : never;

export function karaFuseSearch(
  karas: components["schemas"]["KaraInfoDB"][],
  audioTagMap: TagMap<components["schemas"]["AudioTag"]>,
  videoTagMap: TagMap<components["schemas"]["VideoTag"]>,
  search: string,
) {
  const fuse = new Fuse(karas, {
    keys: [
      { name: "Artists.Name", weight: 2 },
      "Artists.AdditionalNames.Name",
      {
        name: "audioTagsName",
        getFn: (kara) =>
          kara.AudioTags?.map((tag) => audioTagMap.get(tag.ID)?.Name ?? "") ??
          [],
      },
      { name: "Authors.Name", weight: 2 },
      "Comment",
      "ExtraTitles.Name",
      "Language",
      { name: "Medias.name", weight: 2 },
      "Medias.additional_name.Name",
      { name: "SourceMedia.name", weight: 2 },
      { name: "Title", weight: 4 },
      "Version",
      {
        name: "videoTagsName",
        getFn: (kara) =>
          kara.VideoTags?.map((tag) => videoTagMap.get(tag.ID)?.Name ?? "") ??
          [],
      },
    ],
  });
  return fuse.search(search);
}
