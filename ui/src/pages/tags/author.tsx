import { HiSolidTrash } from "solid-icons/hi";
import { createResource, Index } from "solid-js";
import AuthorEditor from "../../components/AuthorEditor";
import { karaberus } from "../../utils/karaberus-client";

export default function TagsAuthor() {
  const [getAuthors, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/author");
    return resp.data;
  });

  const deleteAuthor = async (id: number) => {
    const resp = await karaberus.DELETE("/api/tags/author/{id}", {
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
      <h1 class="header">Author Tags</h1>

      <h2 class="text-2xl font-semibold">Add author</h2>

      <AuthorEditor
        onAdd={() => {
          alert("Author added!");
          refetch();
        }}
      />

      <h2 class="text-2xl font-semibold mt-8">Browse</h2>

      <table class="table">
        <thead>
          <tr>
            <th></th>
            <th>Name</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <Index each={getAuthors()}>
            {(getAuthor) => (
              <tr class="hover">
                <th>{getAuthor().ID}</th>
                <td>{getAuthor().Name}</td>
                <td>
                  <button
                    onclick={() => deleteAuthor(getAuthor().ID)}
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
