import { HiSolidXMark } from "solid-icons/hi";
import {
  createSignal,
  For,
  Show,
  type Accessor,
  type JSX,
  type Setter,
} from "solid-js";

export default function Autocomplete<T>({
  items,
  getItemName,
  getState,
  setState,
  ...inputProps
}: {
  items: T[];
  getItemName: (item: T) => string;
  getState: Accessor<T | undefined>;
  setState: Setter<T | undefined>;
} & JSX.InputHTMLAttributes<HTMLInputElement>) {
  const [getInput, setInput] = createSignal("");

  const filteredItems = () =>
    items.filter((item) =>
      getItemName(item).toLowerCase().includes(getInput().toLowerCase()),
    );

  return (
    <div class="dropdown w-full">
      <div class="textarea textarea-bordered w-full">
        <Show
          when={getState()}
          fallback={
            <input
              type="text"
              value={getInput()}
              oninput={(e) => setInput(e.currentTarget.value)}
              class="outline-none w-full"
              {...inputProps}
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
        <div class="dropdown-content menu bg-base-100 rounded-box z-[1] w-full p-2 shadow max-h-48 overflow-y-auto">
          <ul>
            <For each={filteredItems()}>
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
        </div>
      </Show>
    </div>
  );
}
