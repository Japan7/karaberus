import { HiSolidPencil, HiSolidTrash } from "solid-icons/hi";
import { createResource, Index, useContext } from "solid-js";
import { Context } from "../../components/context";
import MediaEditor from "../../components/MediaEditor";
import type { components } from "../../utils/karaberus";
import { karaberus } from "../../utils/karaberus-client";
import { isAdmin } from "../../utils/session";

export default function TagsMedia() {
  const { getModalRef, setModal, showToast } = useContext(Context)!;

  const [getMedias, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/media");
    return resp.data;
  });

  const postMedia = async (media: components["schemas"]["MediaInfo"]) => {
    const resp = await karaberus.POST("/api/tags/media", { body: media });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    showToast("Media added!");
    refetch();
  };

  const patchMedia =
    (id: components["schemas"]["MediaDB"]["ID"]) =>
    async (media: components["schemas"]["MediaInfo"]) => {
      const resp = await karaberus.PATCH("/api/tags/media/{id}", {
        params: { path: { id } },
        body: media,
      });
      if (resp.error) {
        alert(resp.error.detail);
        return;
      }
      getModalRef().close();
      showToast("Media edited!");
      refetch();
    };

  const deleteMedia = async (id: components["schemas"]["MediaDB"]["ID"]) => {
    if (!confirm("Confirm deletion?")) {
      return;
    }
    const resp = await karaberus.DELETE("/api/tags/media/{id}", {
      params: { path: { id } },
    });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    showToast("Media deleted!");
    refetch();
  };

  return (
    <>
      <h1 class="header">Media Tags</h1>

      <h2 class="text-2xl font-semibold">Add media</h2>

      <MediaEditor onSubmit={postMedia} reset />

      <h2 class="text-2xl font-semibold mt-8">Browse</h2>

      <div class="overflow-auto">
        <table class="table">
          <thead>
            <tr>
              <th></th>
              <th>Type</th>
              <th>Name</th>
              <th>Additional Names</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <Index each={getMedias()}>
              {(getMedia) => (
                <tr class="hover">
                  <th>{getMedia().ID}</th>
                  <td>{getMedia().media_type}</td>
                  <td>{getMedia().name}</td>
                  <td>
                    <ul>
                      <Index each={getMedia().additional_name}>
                        {(getAdditionalName) => (
                          <li>{getAdditionalName().Name}</li>
                        )}
                      </Index>
                    </ul>
                  </td>
                  <td class="flex gap-x-1">
                    <button
                      class="btn btn-sm btn-warning"
                      onclick={() => {
                        setModal(
                          <MediaEditor
                            media={getMedia()}
                            onSubmit={patchMedia(getMedia().ID)}
                          />,
                        );
                        getModalRef().showModal();
                      }}
                    >
                      <HiSolidPencil class="size-4" />
                    </button>
                    <button
                      disabled={!isAdmin()}
                      onclick={() => deleteMedia(getMedia().ID)}
                      class="btn btn-sm btn-error"
                    >
                      <HiSolidTrash class="size-4" />
                    </button>
                  </td>
                </tr>
              )}
            </Index>
          </tbody>
        </table>
      </div>
    </>
  );
}
