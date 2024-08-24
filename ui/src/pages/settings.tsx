import { isTauri } from "@tauri-apps/api/core";
import { HiSolidTrash } from "solid-icons/hi";
import { createResource, createSignal, Index, Show } from "solid-js";
import TokenForm from "../components/TokenForm";
import type { paths } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";
import { getTauriStore, PLAYER_TOKEN_KEY } from "../utils/tauri";

export default function Settings() {
  const [getAllTokens, { refetch: refetchTokens }] = createResource(
    async () => {
      const resp = await karaberus.GET("/api/token");
      return resp.data;
    },
  );
  const [getToken, setToken] = createSignal<string>();

  const deleteToken = async (
    id: paths["/api/token/{token}"]["delete"]["parameters"]["path"]["token"],
  ) => {
    if (!confirm("Confirm deletion?")) {
      return;
    }

    const resp = await karaberus.DELETE("/api/token/{token}", {
      params: {
        path: { token: id },
      },
    });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    refetchTokens();
  };

  const resetPlayerToken = async () => {
    await getTauriStore().delete(PLAYER_TOKEN_KEY);
    await getTauriStore().save();
    await Promise.all(getAllTokens()?.map((t) => deleteToken(t.id)) ?? []);
    refetchTokens();
  };

  return (
    <>
      <h1 class="header">Settings</h1>

      <Show when={isTauri()}>
        <h2 class="text-2xl font-semibold">Reset Player Token</h2>
        <button class="btn btn-warning mb-8" onclick={resetPlayerToken}>
          Reset
        </button>
      </Show>

      <h2 class="text-2xl font-semibold">Create Token</h2>

      <TokenForm
        onToken={(t) => {
          setToken(t);
          refetchTokens();
        }}
      />
      <Show when={getToken()}>
        {(getToken) => (
          <textarea class="textarea textarea-bordered" readonly>
            {getToken()}
          </textarea>
        )}
      </Show>
      <table class="table">
        <thead>
          <tr>
            <th>Token</th>
            <th>Scopes</th>
            <th>Created At</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <Index each={getAllTokens()}>
            {(getToken) => (
              <tr class="hover">
                <th class="font-mono">{getToken().name}</th>
                <td>
                  <ul class="w-full">
                    <li classList={{ hidden: !getToken().scopes.kara }}>
                      kara
                    </li>
                    <li classList={{ hidden: !getToken().scopes.kara_ro }}>
                      kara_ro
                    </li>
                    <li classList={{ hidden: !getToken().scopes.user }}>
                      user
                    </li>
                  </ul>
                </td>
                <td>{new Date(getToken().created_at).toLocaleString()}</td>
                <td>
                  <button
                    onclick={() => deleteToken(getToken().id)}
                    class="btn btn-error btn-sm"
                  >
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
