import { useParams } from "@solidjs/router";
import { createResource, Show } from "solid-js";
import FileUploader from "../../../components/FileUploader";
import KaraEditor from "../../../components/KaraEditor";
import KaraPlayer from "../../../components/KaraPlayer";
import type { components } from "../../../utils/karaberus";
import { karaberus } from "../../../utils/karaberus-client";

export default function KaraokeBrowseId() {
  const params = useParams();

  const [getKara] = createResource(async () => {
    const resp = await karaberus.GET("/api/kara/{id}", {
      params: {
        path: {
          id: parseInt(params.id),
        },
      },
    });
    return resp.data;
  });

  const onUpload = () => {
    alert("Upload complete!");
    location.reload();
  };

  const edit = async (info: components["schemas"]["KaraInfo"]) => {
    const resp = await karaberus.PATCH("/api/kara/{id}", {
      body: info,
      params: {
        path: {
          id: parseInt(params.id),
        },
      },
    });

    if (resp.error) {
      alert(resp.error);
      return;
    }

    location.reload();
  };

  return (
    <>
      <KaraPlayer id={params.id} />

      <h2 class="text-2xl font-semibold mt-4">Upload files</h2>

      <div class="grid md:grid-cols-3 gap-x-2">
        <FileUploader
          title="Video"
          method="PUT"
          url={`/api/kara/${params.id}/upload/video`}
          downloadUrl={`/api/kara/${params.id}/download/video`}
          onUpload={onUpload}
        />
        <FileUploader
          title="Instrumental"
          method="PUT"
          url={`/api/kara/${params.id}/upload/inst`}
          downloadUrl={`/api/kara/${params.id}/download/inst`}
          onUpload={onUpload}
        />
        <FileUploader
          title="Subtitles"
          method="PUT"
          url={`/api/kara/${params.id}/upload/sub`}
          downloadUrl={`/api/kara/${params.id}/download/sub`}
          onUpload={onUpload}
        />
      </div>

      <h2 class="text-2xl font-semibold mt-4">Edit</h2>

      <Show when={getKara()?.kara} fallback={<p>loading karaoke...</p>}>
        {(getKara) => <KaraEditor kara={getKara()} onSubmit={edit} />}
      </Show>
    </>
  );
}
