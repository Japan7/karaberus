import { HiSolidTrash } from "solid-icons/hi";
import { createResource, createSignal, Index, Show } from "solid-js";
import MediaEditor from "../../components/MediaEditor";
import { karaberus } from "../../utils/karaberus-client";
import { isAdmin } from "../../utils/session";

export default function TagsMedia() {
  const [getMedias, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/media");
    return resp.data;
  });

  const [getToast, setToast] = createSignal<string>();

  const deleteArtist = async (id: number) => {
    const resp = await karaberus.DELETE("/api/tags/media/{id}", {
      params: { path: { id } },
    });
    if (resp.error) {
      alert(resp.error.title);
      return;
    }
    refetch();
  };

  return (
    <>
      <h1 class="header">Media Tags</h1>

      <h2 class="text-2xl font-semibold">Add media</h2>

      <MediaEditor
        onAdd={() => {
          setToast("Media added!");
          setTimeout(() => setToast(), 3000);
          refetch();
        }}
      />

      <h2 class="text-2xl font-semibold mt-8">Browse</h2>

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
                <td>
                  <button
                    disabled={!isAdmin()}
                    onclick={() => deleteArtist(getMedia().ID)}
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
