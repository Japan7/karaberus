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
      "Artists.Name",
      "Artists.AdditionalNames.Name",
      {
        name: "audioTagsName",
        getFn: (kara) =>
          kara.AudioTags?.map((tag) => audioTagMap.get(tag.ID)?.Name ?? "") ??
          [],
      },
      "Authors.Name",
      "Comment",
      "ExtraTitles.Name",
      "Language",
      "Medias.name",
      "Medias.additional_name.Name",
      "SourceMedia.name",
      "Title",
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
