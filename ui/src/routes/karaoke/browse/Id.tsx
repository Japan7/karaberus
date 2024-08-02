import { useParams } from "@solidjs/router";
import KaraPlayer from "../../../components/KaraPlayer";
import { fileForm, karaberus } from "../../../utils/karaberus-client";

export default function KaraokeBrowseId() {
  const params = useParams();

  const upload = async (filetype: string, file?: File) => {
    if (!file) return;

    const resp = await karaberus.PUT("/api/kara/{id}/upload/{filetype}", {
      params: {
        path: {
          id: parseInt(params.id),
          filetype,
        },
      },
      ...fileForm(file),
    });

    if (resp.error) {
      alert(resp.error);
      return;
    }

    alert(filetype + " uploaded successfully");
    location.reload();
  };

  return (
    <>
      <input
        type="file"
        onchange={(e) => upload("video", e.target.files?.[0])}
        class="file-input file-input-bordered w-full max-w-xs"
      />
      <input
        type="file"
        onchange={(e) => upload("sub", e.target.files?.[0])}
        class="file-input file-input-bordered w-full max-w-xs"
      />

      <KaraPlayer id={params.id} />
    </>
  );
}
