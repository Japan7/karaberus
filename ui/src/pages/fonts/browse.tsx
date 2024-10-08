import { HiSolidTrash } from "solid-icons/hi";
import { createResource, Index, useContext } from "solid-js";
import { Context } from "../../components/Context";
import DownloadAnchor from "../../components/DownloadAnchor";
import FileUploader from "../../components/FileUploader";
import { apiUrl, karaberus } from "../../utils/karaberus-client";

export default function FontsBrowse() {
  const { showToast } = useContext(Context)!;

  const [getFonts, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/font");
    return resp.data?.Fonts;
  });

  const onUpload = () => {
    showToast("Font uploaded!");
    refetch();
  };

  return (
    <>
      <h1 class="header">Fonts</h1>

      <h2 class="text-2xl font-semibold">Upload</h2>

      <FileUploader
        method="POST"
        url={apiUrl("api/font")}
        onUpload={onUpload}
      />

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
                  <DownloadAnchor
                    href={apiUrl(`api/font/${getFont().ID}/download`)}
                    download={getFont().Name}
                  >
                    {getFont().Name}
                  </DownloadAnchor>
                </td>
                <td>{new Date(getFont().UpdatedAt).toLocaleString()}</td>
                <td>
                  <button disabled class="btn btn-error btn-sm">
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
