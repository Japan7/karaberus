import { useParams } from "@solidjs/router";
import { createResource, Show } from "solid-js";
import FileUploader from "../../../components/FileUploader";
import KaraEditor from "../../../components/KaraEditor";
import KaraPlayer from "../../../components/KaraPlayer";
import type { components } from "../../../utils/karaberus";
import { karaberus } from "../../../utils/karaberus-client";
import { isAdmin } from "../../../utils/session";
import { HiSolidTrash } from "solid-icons/hi";

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

  const editKaraoke = async (info: components["schemas"]["KaraInfo"]) => {
    const resp = await karaberus.PATCH("/api/kara/{id}", {
      body: info,
      params: {
        path: {
          id: parseInt(params.id),
        },
      },
    });

    if (resp.error) {
      alert(resp.error.title);
      return;
    }

    location.reload();
  };

  const deleteKaraoke = async () => {
    const resp = await karaberus.DELETE("/api/kara/{id}", {
      params: {
        path: {
          id: parseInt(params.id),
        },
      },
    });

    if (resp.error) {
      alert(resp.error.title);
      return;
    }

    location.href = "/karaoke/browse";
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

      <h2 class="text-2xl font-semibold mt-4">Edit metadata</h2>

      <Show when={getKara()?.kara} fallback={<p>loading karaoke...</p>}>
        {(getKara) => <KaraEditor kara={getKara()} onSubmit={editKaraoke} />}
      </Show>

      <div class="divider" />

      <div class="collapse collapse-arrow bg-neutral text-neutral-content">
        <input type="checkbox" />
        <div class="collapse-title text-xl font-medium">Danger zone</div>
        <div class="collapse-content">
          <button
            disabled={!isAdmin()}
            onclick={deleteKaraoke}
            class="btn btn-error w-full"
          >
            <HiSolidTrash class="size-5" />
            Delete Karaoke
          </button>
        </div>
      </div>
    </>
  );
}
