import { useParams } from "@solidjs/router";
import { createResource, Show } from "solid-js";
import KaraEditor from "../../../components/KaraEditor";
import KaraFileUploader from "../../../components/KaraFileUploader";
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
        <KaraFileUploader
          title="Video"
          putUrl={`/api/kara/${params.id}/upload/video`}
          downloadUrl={`/api/kara/${params.id}/download/video`}
        />
        <KaraFileUploader
          title="Instrumental"
          putUrl={`/api/kara/${params.id}/upload/inst`}
          downloadUrl={`/api/kara/${params.id}/download/inst`}
        />
        <KaraFileUploader
          title="Subtitles"
          putUrl={`/api/kara/${params.id}/upload/sub`}
          downloadUrl={`/api/kara/${params.id}/download/sub`}
        />
      </div>

      <h2 class="text-2xl font-semibold mt-4">Edit</h2>

      <Show when={getKara()?.kara} fallback={<p>loading karaoke...</p>}>
        {(getKara) => <KaraEditor kara={getKara()} onSubmit={edit} />}
      </Show>
    </>
  );
}
