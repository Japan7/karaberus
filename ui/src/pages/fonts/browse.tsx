import { HiSolidTrash } from "solid-icons/hi";
import { createResource, Index } from "solid-js";
import FileUploader from "../../components/FileUploader";
import { karaberus } from "../../utils/karaberus-client";

export default function FontsBrowse() {
  const [getFonts, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/font");
    return resp.data?.Fonts;
  });

  const onUpload = () => {
    alert("Font uploaded!");
    refetch();
  };

  return (
    <>
      <h1 class="header">Fonts</h1>

      <h2 class="text-2xl font-semibold">Upload</h2>

      <FileUploader method="POST" url="/api/font" onUpload={onUpload} />

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
                <td>
                  <a
                    href={"/api/font/" + getFont().ID}
                    download={getFont().Name}
                    class="link"
                  >
                    {getFont().Name}
                  </a>
                </td>
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
