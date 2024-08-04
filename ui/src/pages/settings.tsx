import { createResource, createSignal, Index, Show } from "solid-js";
import TokenForm from "../components/TokenForm";
import { karaberus } from "../utils/karaberus-client";
import { HiSolidTrash } from "solid-icons/hi";

export default function Settings() {
  const [getAllTokens, { refetch: refetchTokens }] = createResource(
    async () => {
      const resp = await karaberus.GET("/api/token");
      return resp.data;
    },
  );
  const [getToken, setToken] = createSignal<string>();

  const deleteToken = async (id: string) => {
    const resp = await karaberus.DELETE("/api/token/{token}", {
      params: {
        path: { token: id },
      },
    });
    if (resp.error) {
      alert(resp.error);
      return;
    }
    refetchTokens();
  };

  return (
    <>
      <h1 class="header">Settings</h1>

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
                <th class="font-mono break-all max-w-xs">{getToken().ID}</th>
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
                <td>{new Date(getToken().CreatedAt).toLocaleString()}</td>
                <td>
                  <button
                    onclick={() => deleteToken(getToken().ID)}
                    class="btn btn-sm"
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
