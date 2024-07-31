import { HiSolidXMark } from "solid-icons/hi";
import { createSignal, For, Show, type Accessor, type Setter } from "solid-js";

export default function Autocomplete<T>({
  items,
  getItemName,
  getState,
  setState,
  required,
}: {
  items: T[];
  getItemName: (item: T) => string;
  getState: Accessor<T | undefined>;
  setState: Setter<T | undefined>;
  required?: boolean;
}) {
  const [getInput, setInput] = createSignal("");

  return (
    <div class="dropdown w-full">
      <div class="textarea textarea-bordered w-full">
        <Show
          when={getState()}
          fallback={
            <input
              type="text"
              required={required}
              value={getInput()}
              oninput={(e) => setInput(e.currentTarget.value)}
              class="outline-none w-full"
            />
          }
        >
          {(getState) => (
            <div
              onclick={() => setState(undefined)}
              class="btn w-full flex justify-evenly"
            >
              <span>{getItemName(getState())}</span>
              <HiSolidXMark class="size-5" />
            </div>
          )}
        </Show>
      </div>
      <Show when={getState() === undefined}>
        <ul class="dropdown-content menu bg-base-100 rounded-box z-[1] w-full p-2 shadow">
          <For each={items}>
            {(item) => (
              <li>
                <a
                  href="#"
                  onclick={(e) => {
                    e.preventDefault();
                    setState(() => item);
                    setInput("");
                  }}
                >
                  {getItemName(item)}
                </a>
              </li>
            )}
          </For>
        </ul>
      </Show>
    </div>
  );
}
