import { createResource, Show, useContext } from "solid-js";
import { Context } from "../components/Context";
import ProfileEditor from "../components/ProfileEditor";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";
import { getSessionInfos } from "../utils/session";

export default function Profile() {
  const { showToast } = useContext(Context)!;

  const [getMe, { refetch: refetchMe }] = createResource(async () => {
    const resp = await karaberus.GET("/api/me");
    return resp.data;
  });

  const infos = getSessionInfos();

  const onSubmit = async (data: {
    author?: components["schemas"]["TimingAuthor"];
  }) => {
    const resp = await karaberus.PUT("/api/me/author", {
      body: {
        id: data.author?.ID || null,
      },
    });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    showToast("Your profile is updated!");
    refetchMe();
  };

  return (
    <>
      <h1 class="header">Profile</h1>

      <Show
        when={getMe()}
        fallback={<span class="loading loading-spinner loading-lg" />}
      >
        {(getMe) => <ProfileEditor user={getMe()} onSubmit={onSubmit} />}
      </Show>

      <h2 class="text-2xl font-semibold mt-4">User</h2>
      <pre>{JSON.stringify(getMe(), undefined, 2)}</pre>

      <h2 class="text-2xl font-semibold mt-4">JWT</h2>
      <pre>{JSON.stringify(infos, undefined, 2)}</pre>
    </>
  );
}
