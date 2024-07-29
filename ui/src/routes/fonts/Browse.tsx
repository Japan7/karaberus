import { HiSolidTrash } from "solid-icons/hi";
import { createResource, Index } from "solid-js";
import { karaberus } from "../../utils/karaberus-client";

export default function FontsBrowse() {
  const [getFonts] = createResource(async () => {
    const resp = await karaberus.GET("/api/font");
    return resp.data?.Fonts;
  });

  return (
    <>
      <h1 class="header">Browse Fonts</h1>

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
            {(getItem) => (
              <tr class="hover">
                <th>{getItem().ID}</th>
                <td>{getItem().Name}</td>
                <td>{new Date(getItem().UpdatedAt).toLocaleString()}</td>
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
