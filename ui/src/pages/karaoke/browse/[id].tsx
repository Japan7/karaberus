import { useNavigate, useParams } from "@solidjs/router";
import { isTauri } from "@tauri-apps/api/core";
import {
  HiOutlineArrowLeft,
  HiOutlineArrowTopRightOnSquare,
  HiOutlineCheck,
  HiOutlineLink,
  HiSolidTrash,
} from "solid-icons/hi";
import { createResource, createSignal, Show } from "solid-js";
import BrowserKaraPlayer from "../../../components/BrowserKaraPlayer";
import DownloadAnchor from "../../../components/DownloadAnchor";
import FileUploader from "../../../components/FileUploader";
import KaraEditor from "../../../components/KaraEditor";
import MpvKaraPlayer from "../../../components/MpvKaraPlayer";
import type { components } from "../../../utils/karaberus";
import { apiUrl, karaberus } from "../../../utils/karaberus-client";
import { isAdmin } from "../../../utils/session";
import { buildKaraberusUrl, IS_TAURI_DIST_BUILD } from "../../../utils/tauri";

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

  const [getCopySuccess, setCopySuccess] = createSignal(false);

  const copyLink = () => {
    const link = IS_TAURI_DIST_BUILD
      ? buildKaraberusUrl(location.pathname).toString()
      : location.href;
    navigator.clipboard.writeText(link);
    setCopySuccess(true);
    setTimeout(() => setCopySuccess(false), 1000);
  };

  return (
    <>
      <div class="flex justify-between">
        <button onclick={() => navigate(-1)} class="btn btn-sm w-fit btn-ghost">
          <HiOutlineArrowLeft class="size-5" />
          Back
        </button>

        <button
          onclick={copyLink}
          class="btn btn-sm w-fit"
          classList={{ "btn-success": getCopySuccess() }}
        >
          {getCopySuccess() ? (
            <HiOutlineCheck class="size-5" />
          ) : (
            <HiOutlineLink class="size-5" />
          )}
          Copy link
        </button>
      </div>

      <Show when={getKara()}>
        {(getKara) =>
          isTauri() ? (
            <MpvKaraPlayer kara={getKara().kara} />
          ) : (
            <BrowserKaraPlayer kara={getKara().kara} />
          )
        }
      </Show>

      <h2 class="text-2xl font-semibold mt-4">Upload files</h2>

      <div class="flex flex-col md:grid md:grid-cols-3 md:gap-x-2">
        <FileUploader
          title="Video"
          method="PUT"
          url={apiUrl(`api/kara/${params.id}/upload/video`)}
          onUpload={onUpload}
          altChildren={
            getKara()?.kara.VideoUploaded ? (
              <DownloadAnchor
                href={apiUrl(`api/kara/${params.id}/download/video`)}
                download={`${params.id}.mkv`}
              >
                {"Download current video "}
                <HiOutlineArrowTopRightOnSquare class="size-4" />
              </DownloadAnchor>
            ) : (
              "(No upload yet)"
            )
          }
        />
        <FileUploader
          title="Instrumental"
          method="PUT"
          url={apiUrl(`api/kara/${params.id}/upload/inst`)}
          onUpload={onUpload}
          altChildren={
            getKara()?.kara.InstrumentalUploaded ? (
              <DownloadAnchor
                href={apiUrl(`api/kara/${params.id}/download/inst`)}
                download={`${params.id}.mka`}
              >
                {"Download current instrumental "}
                <HiOutlineArrowTopRightOnSquare class="size-4" />
              </DownloadAnchor>
            ) : (
              "(No upload yet)"
            )
          }
        />
        <FileUploader
          title="Subtitles"
          method="PUT"
          url={apiUrl(`api/kara/${params.id}/upload/sub`)}
          onUpload={onUpload}
          altChildren={
            getKara()?.kara.SubtitlesUploaded ? (
              <DownloadAnchor
                href={apiUrl(`api/kara/${params.id}/download/sub`)}
                download={`${params.id}.ass`}
              >
                {"Download current subtitles "}
                <HiOutlineArrowTopRightOnSquare class="size-4" />
              </DownloadAnchor>
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
