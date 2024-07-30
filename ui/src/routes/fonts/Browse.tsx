import { HiSolidTrash } from "solid-icons/hi";
import { createResource, createSignal, Index, type JSX } from "solid-js";
import { fileForm, karaberus } from "../../utils/karaberus-client";

export default function FontsBrowse() {
  const [getFonts, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/font");
    return resp.data?.Fonts;
  });

  const [getFile, setFile] = createSignal<File>();

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();
    const file = getFile();
    if (!file) {
      alert("Please select a file");
      return;
    }
    const resp = await karaberus.POST("/api/font", fileForm(file));
    if (resp.error) {
      alert(resp.error);
      return;
    }
    setFile(undefined);
    refetch();
  };

  return (
    <>
      <h1 class="header">Fonts</h1>

      <h2 class="text-2xl font-semibold">Upload</h2>

      <form onsubmit={onsubmit} class="flex gap-x-2 mt-4">
        <input
          type="file"
          required
          onchange={(e) => setFile(e.target.files?.[0])}
          class="file-input file-input-bordered w-full max-w-xs"
        />
        <input type="submit" class="btn" />
      </form>

      <h2 class="text-2xl font-semibold mt-8">Browse</h2>

      <table class="table">
        <thead>
          <tr>
            <th></th>
            <th>Name</th>
            <th>Updated At</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <Index each={getFonts()}>
            {(getFont) => (
              <tr class="hover">
                <th>{getFont().ID}</th>
                <td>{getFont().Name}</td>
                <td>{new Date(getFont().UpdatedAt).toLocaleString()}</td>
                <td>
                  <button disabled class="btn btn-sm">
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
