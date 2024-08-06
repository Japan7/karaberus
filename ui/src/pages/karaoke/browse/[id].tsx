import { useNavigate, useParams } from "@solidjs/router";
import { HiOutlineArrowTopRightOnSquare, HiSolidTrash } from "solid-icons/hi";
import { createResource, Show } from "solid-js";
import FileUploader from "../../../components/FileUploader";
import KaraEditor from "../../../components/KaraEditor";
import KaraPlayer from "../../../components/KaraPlayer";
import type { components } from "../../../utils/karaberus";
import { karaberus } from "../../../utils/karaberus-client";
import { isAdmin } from "../../../utils/session";

export default function KaraokeBrowseId() {
  const params = useParams();
  const navigate = useNavigate();

  const [getKara, { refetch }] = createResource(async () => {
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
    refetch();
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

    alert("Karaoke updated!");
    refetch();
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

    navigate("/karaoke/browse");
  };

  return (
    <>
      <Show when={getKara()}>
        {(getKara) => <KaraPlayer kara={getKara().kara} />}
      </Show>

      <h2 class="text-2xl font-semibold mt-4">Upload files</h2>

      <div class="flex flex-col md:grid md:grid-cols-3 md:gap-x-2">
        <FileUploader
          title="Video"
          method="PUT"
          url={`/api/kara/${params.id}/upload/video`}
          onUpload={onUpload}
          altChildren={
            getKara()?.kara.VideoUploaded ? (
              <a
                href={`/api/kara/${params.id}/download/video`}
                target="_blank"
                rel="noreferrer"
                class="link flex gap-x-1"
              >
                Open current <HiOutlineArrowTopRightOnSquare class="size-4" />
              </a>
            ) : (
              "(No upload yet)"
            )
          }
        />
        <FileUploader
          title="Instrumental"
          method="PUT"
          url={`/api/kara/${params.id}/upload/inst`}
          onUpload={onUpload}
          altChildren={
            getKara()?.kara.InstrumentalUploaded ? (
              <a
                href={`/api/kara/${params.id}/download/inst`}
                target="_blank"
                rel="noreferrer"
                class="link flex gap-x-1"
              >
                Open current <HiOutlineArrowTopRightOnSquare class="size-4" />
              </a>
            ) : (
              "(No upload yet)"
            )
          }
        />
        <FileUploader
          title="Subtitles"
          method="PUT"
          url={`/api/kara/${params.id}/upload/sub`}
          onUpload={onUpload}
          altChildren={
            getKara()?.kara.SubtitlesUploaded ? (
              <a
                href={`/api/kara/${params.id}/download/sub`}
                target="_blank"
                rel="noreferrer"
                class="link flex gap-x-1"
              >
                Open current <HiOutlineArrowTopRightOnSquare class="size-4" />
              </a>
            ) : (
              "(No upload yet)"
            )
          }
        />
      </div>

      <h2 class="text-2xl font-semibold mt-4">Edit metadata</h2>

      <Show
        when={getKara()?.kara}
        fallback={<span class="loading loading-spinner loading-lg" />}
      >
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
