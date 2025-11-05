import { createSignal, type JSX } from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";

export default function TokenForm(props: {
  onToken: (
    token: components["schemas"]["CreateTokenOutputBody"]["token"],
  ) => void;
}) {
  const [getName, setName] = createSignal("");
  const [getKaraChecked, setKaraChecked] = createSignal(false);
  const [getKaraROChecked, setKaraROChecked] = createSignal(false);
  const [getUserChecked, setUserChecked] = createSignal(false);

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();

    const resp = await karaberus.POST("/api/token", {
      body: {
        name: getName(),
        scopes: {
          kara: getKaraChecked(),
          kara_ro: getKaraROChecked(),
          user: getUserChecked(),
        },
      },
    });

    if (resp.error) {
      alert(resp.error.detail);
      return;
    }

    setName("");
    setKaraChecked(false);
    setKaraROChecked(false);
    setUserChecked(false);
    props.onToken(resp.data.token);
  };

  return (
    <form onSubmit={onsubmit} class="flex flex-col gap-y-1 w-full">
      <input
        type="text"
        required
        placeholder="Name"
        value={getName()}
        oninput={(e) => setName(e.currentTarget.value)}
        class="input"
      />
      <label class="label cursor-pointer">
        <span>Karaoke</span>
        <input
          type="checkbox"
          checked={getKaraChecked()}
          onchange={(e) => setKaraChecked(e.currentTarget.checked)}
          class="toggle"
        />
      </label>
      <label class="label cursor-pointer">
        <span>Karaoke Read-Only</span>
        <input
          type="checkbox"
          checked={getKaraROChecked()}
          onchange={(e) => setKaraROChecked(e.currentTarget.checked)}
          class="toggle"
        />
      </label>
      <label class="label cursor-pointer">
        <span>User</span>
        <input
          type="checkbox"
          checked={getUserChecked()}
          onchange={(e) => setUserChecked(e.currentTarget.checked)}
          class="toggle"
        />
      </label>
      <input type="submit" class="btn btn-primary w-full" />
    </form>
  );
}
