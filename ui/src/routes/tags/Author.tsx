import { HiSolidTrash } from "solid-icons/hi";
import { createResource, createSignal, Index, type JSX } from "solid-js";
import { karaberus } from "../../utils/karaberus-client";

export default function TagsAuthor() {
  const [getAuthors, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/author");
    return resp.data;
  });

  const [getName, setName] = createSignal("");

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();
    const resp = await karaberus.POST("/api/tags/author", {
      body: { name: getName() },
    });
    if (resp.error) {
      alert(resp.error);
      return;
    }
    setName("");
    refetch();
  };

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

      <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full max-w-xs">
        <label>
          <div class="label">
            <span class="label-text">Name</span>
          </div>
          <input
            type="text"
            required
            value={getName()}
            onInput={(e) => setName(e.currentTarget.value)}
            class="input input-bordered w-full"
          />
        </label>

        <input type="submit" class="btn" />
      </form>

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
