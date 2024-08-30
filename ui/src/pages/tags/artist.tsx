import { HiSolidPencil, HiSolidTrash } from "solid-icons/hi";
import { createResource, Index, useContext } from "solid-js";
import ArtistEditor from "../../components/ArtistEditor";
import { Context } from "../../components/Context";
import type { components } from "../../utils/karaberus";
import { karaberus } from "../../utils/karaberus-client";
import { isAdmin } from "../../utils/session";

export default function TagsArtist() {
  const { getModalRef, setModal, showToast } = useContext(Context)!;

  const [getArtists, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/artist");
    return resp.data;
  });

  const postArtist = async (artist: components["schemas"]["ArtistInfo"]) => {
    const resp = await karaberus.POST("/api/tags/artist", { body: artist });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    showToast("Artist added!");
    refetch();
  };

  const patchArtist =
    (id: components["schemas"]["Artist"]["ID"]) =>
    async (artist: components["schemas"]["ArtistInfo"]) => {
      const resp = await karaberus.PATCH("/api/tags/artist/{id}", {
        params: { path: { id } },
        body: artist,
      });
      if (resp.error) {
        alert(resp.error.detail);
        return;
      }
      getModalRef().close();
      showToast("Artist edited!");
      refetch();
    };

  const deleteArtist = async (id: components["schemas"]["Artist"]["ID"]) => {
    if (!confirm("Confirm deletion?")) {
      return;
    }
    const resp = await karaberus.DELETE("/api/tags/artist/{id}", {
      params: { path: { id } },
    });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    showToast("Artist deleted!");
    refetch();
  };

  return (
    <>
      <h1 class="header">Artist Tags</h1>

      <h2 class="text-2xl font-semibold">Add artist</h2>

      <ArtistEditor onSubmit={postArtist} reset />

      <h2 class="text-2xl font-semibold mt-8">Browse</h2>

      <div class="overflow-auto">
        <table class="table">
          <thead>
            <tr>
              <th></th>
              <th>Name</th>
              <th>Additional Names</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <Index each={getArtists()}>
              {(getArtist) => (
                <tr class="hover">
                  <th>{getArtist().ID}</th>
                  <td>{getArtist().Name}</td>
                  <td>
                    <ul>
                      <Index each={getArtist().AdditionalNames}>
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
                          <ArtistEditor
                            artist={getArtist()}
                            onSubmit={patchArtist(getArtist().ID)}
                          />,
                        );
                        getModalRef().showModal();
                      }}
                    >
                      <HiSolidPencil class="size-4" />
                    </button>
                    <button
                      disabled={!isAdmin()}
                      onclick={() => deleteArtist(getArtist().ID)}
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
