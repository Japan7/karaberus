import { Icon } from "solid-heroicons";
import { trash } from "solid-heroicons/solid";
import { createResource } from "solid-js";
import { karaberus } from "../../utils/karaberus-client";

export default function FontsBrowse() {
  const [fonts, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/font");
    return resp.data?.Fonts;
  });

  return (
    <>
      <h1 class="text-6xl font-bold mt-16 mb-8">Browse Fonts</h1>

      <table class="table table-zebra">
        <thead>
          <tr>
            <th></th>
            <th>Name</th>
            <th>Uploaded At</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {fonts()?.map((font) => (
            <tr class="hover">
              <th>{font.ID}</th>
              <td>{font.Name}</td>
              <td>{font.UploadedAt}</td>
              <td>
                <button disabled class="btn btn-sm">
                  <Icon path={trash} class="size-4" />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}
