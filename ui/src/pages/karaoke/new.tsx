import { useNavigate } from "@solidjs/router";
import KaraEditor from "../../components/KaraEditor";
import type { components } from "../../utils/karaberus";
import { karaberus } from "../../utils/karaberus-client";
import { createResource, Show } from "solid-js";

export default function KaraokeNew() {
  const navigate = useNavigate();

  const [getMe] = createResource(async () => {
    const resp = await karaberus.GET("/api/me");
    return resp.data;
  });

  const onSubmit = async (info: components["schemas"]["KaraInfo"]) => {
    const resp = await karaberus.POST("/api/kara", {
      body: info,
    });

    if (resp.error) {
      alert(resp.error.detail);
      return;
    }

    navigate("/karaoke/browse/" + resp.data.kara.ID);
  };

  return (
    <>
      <h1 class="header">New Karaoke</h1>

      <Show
        when={getMe()}
        fallback={<span class="loading loading-spinner loading-lg" />}
      >
        {(getMe) => <KaraEditor onSubmit={onSubmit} me={getMe()} />}
      </Show>
    </>
  );
}
