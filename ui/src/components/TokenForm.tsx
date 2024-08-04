import { createSignal, type JSX } from "solid-js";
import { karaberus } from "../utils/karaberus-client";

export default function TokenForm(props: { onToken: (token: string) => void }) {
  const [getKaraChecked, setKaraChecked] = createSignal(false);
  const [getKaraROChecked, setKaraROChecked] = createSignal(false);
  const [getUserChecked, setUserChecked] = createSignal(false);

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = async (e) => {
    e.preventDefault();

    const resp = await karaberus.POST("/api/token", {
      body: {
        kara: getKaraChecked(),
        kara_ro: getKaraROChecked(),
        user: getUserChecked(),
      },
    });

    if (resp.error) {
      alert(resp.error);
      return;
    }

    setKaraChecked(false);
    setKaraROChecked(false);
    setUserChecked(false);
    props.onToken(resp.data.token);
  };

  return (
    <form onSubmit={onsubmit} class="flex flex-col gap-y-1 max-w-xs">
      <div class="form-control">
        <label class="label cursor-pointer">
          <span class="label-text">Karaoke</span>
          <input
            type="checkbox"
            checked={getKaraChecked()}
            onchange={(e) => setKaraChecked(e.currentTarget.checked)}
            class="toggle"
          />
        </label>
      </div>
      <div class="form-control">
        <label class="label cursor-pointer">
          <span class="label-text">Karaoke Read-Only</span>
          <input
            type="checkbox"
            checked={getKaraROChecked()}
            onchange={(e) => setKaraROChecked(e.currentTarget.checked)}
            class="toggle"
          />
        </label>
      </div>
      <div class="form-control">
        <label class="label cursor-pointer">
          <span class="label-text">User</span>
          <input
            type="checkbox"
            checked={getUserChecked()}
            onchange={(e) => setUserChecked(e.currentTarget.checked)}
            class="toggle"
          />
        </label>
      </div>
      <input type="submit" class="btn w-full" />
    </form>
  );
}
