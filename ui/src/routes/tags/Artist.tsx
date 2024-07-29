import { HiSolidTrash } from "solid-icons/hi";
import { createResource, createSignal, Index, type JSX } from "solid-js";
import { karaberus } from "../../utils/karaberus-client";

export default function TagsArtist() {
  const [getArtists, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/artist");
    return resp.data;
  });

  const [getName, setName] = createSignal("");
  const [getAdditionalNames, setAdditionalNames] = createSignal("");

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();
    const resp = await karaberus.POST("/api/tags/artist", {
      body: {
        name: getName(),
        additional_names: getAdditionalNames().trim().split("\n"),
      },
    });
    if (resp.error) {
      alert(resp.error);
      return;
    }
    setName("");
    setAdditionalNames("");
    refetch();
  };

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

        <label>
          <div class="label">
            <span class="label-text">Additional names</span>
            <span class="label-text-alt">1 per line</span>
          </div>
          <textarea
            value={getAdditionalNames()}
            onInput={(e) => setAdditionalNames(e.currentTarget.value)}
            class="textarea textarea-bordered w-full"
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
                  <ul class="list-disc list-inside">
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
