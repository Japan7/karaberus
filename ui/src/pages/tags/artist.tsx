import { HiSolidTrash } from "solid-icons/hi";
import { createResource, Index } from "solid-js";
import ArtistEditor from "../../components/ArtistEditor";
import { karaberus } from "../../utils/karaberus-client";

export default function TagsArtist() {
  const [getArtists, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/artist");
    return resp.data;
  });

  const deleteArtist = async (id: number) => {
    const resp = await karaberus.DELETE("/api/tags/artist/{id}", {
      params: { path: { id } },
    });
    if (resp.error) {
      alert(resp.error);
      return;
    }
    refetch();
  };

  return (
    <>
      <h1 class="header">Artist Tags</h1>

      <h2 class="text-2xl font-semibold">Add artist</h2>

      <ArtistEditor
        onAdd={() => {
          alert("Artist added!");
          refetch();
        }}
      />

      <h2 class="text-2xl font-semibold mt-8">Browse</h2>

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
                <td>
                  <button
                    onclick={() => deleteArtist(getArtist().ID)}
                    class="btn btn-sm"
                  >
                    <HiSolidTrash class="size-4" />
                  </button>
                </td>
              </tr>
            )}
          </Index>
        </tbody>
      </table>
    </>
  );
}
