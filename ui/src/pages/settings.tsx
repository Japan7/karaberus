import { createSignal, Show } from "solid-js";
import TokenForm from "../components/TokenForm";

export default function Settings() {
  const [getToken, setToken] = createSignal<string>();

  return (
    <>
      <h1 class="header">Settings</h1>

      <h2 class="text-2xl font-semibold">Create Token</h2>

      <TokenForm onToken={setToken} />
      <Show when={getToken()}>
        {(getToken) => (
          <textarea class="textarea textarea-bordered" readonly>
            {getToken()}
          </textarea>
        )}
      </Show>
    </>
  );
}
