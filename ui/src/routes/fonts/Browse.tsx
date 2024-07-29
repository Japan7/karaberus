import { HiSolidTrash } from "solid-icons/hi";
import { createResource } from "solid-js";
import { karaberus } from "../../utils/karaberus-client";

export default function FontsBrowse() {
  const [getFonts] = createResource(async () => {
    const resp = await karaberus.GET("/api/font");
    return resp.data?.Fonts;
  });

  return (
    <>
      <h1 class="text-6xl font-bold mt-16 mb-8">Browse Fonts</h1>

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
          {getFonts()?.map((font) => (
            <tr class="hover">
              <th>{font.ID}</th>
              <td>{font.Name}</td>
              <td>{new Date(font.UpdatedAt).toLocaleString()}</td>
              <td>
                <button disabled class="btn btn-sm">
                  <HiSolidTrash class="size-4" />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}
