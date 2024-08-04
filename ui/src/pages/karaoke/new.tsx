import { useNavigate } from "@solidjs/router";
import KaraEditor from "../../components/KaraEditor";
import type { components } from "../../utils/karaberus";
import { karaberus } from "../../utils/karaberus-client";

export default function KaraokeNew() {
  const navigate = useNavigate();

  const onSubmit = async (info: components["schemas"]["KaraInfo"]) => {
    const resp = await karaberus.POST("/api/kara", {
      body: info,
    });

    if (resp.error) {
      alert(resp.error);
      return;
    }

    navigate("/karaoke/browse/" + resp.data.kara.ID);
  };

  return (
    <>
      <h1 class="header">New Karaoke</h1>

      <KaraEditor onSubmit={onSubmit} />
    </>
  );
}
