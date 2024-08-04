import { HiSolidXMark } from "solid-icons/hi";
import {
  createSignal,
  For,
  Show,
  splitProps,
  type Accessor,
  type JSX,
  type Setter,
} from "solid-js";

export default function Autocomplete<T>(
  props: {
    items: T[];
    getItemName: (item: T) => string;
    getState: Accessor<T | undefined>;
    setState: Setter<T | undefined>;
  } & JSX.InputHTMLAttributes<HTMLInputElement>,
) {
  const [local, inputProps] = splitProps(props, [
    "items",
    "getItemName",
    "getState",
    "setState",
  ]);

  const [getInput, setInput] = createSignal("");

  const filteredItems = () =>
    local.items.filter((item) =>
      local.getItemName(item).toLowerCase().includes(getInput().toLowerCase()),
    );

  return (
    <div class="dropdown w-full">
      <div class="textarea textarea-bordered w-full">
        <Show
          when={local.getState()}
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
              onclick={() => local.setState(undefined)}
              class="btn w-full flex justify-evenly"
            >
              <span>{local.getItemName(getState())}</span>
              <HiSolidXMark class="size-5" />
            </div>
          )}
        </Show>
      </div>
      <Show when={local.getState() === undefined}>
        <div class="dropdown-content menu bg-base-100 rounded-box z-[1] w-full p-2 shadow max-h-48 overflow-y-auto">
          <ul>
            <For each={filteredItems()}>
              {(item) => (
                <li>
                  <a
                    href="#"
                    onclick={(e) => {
                      e.preventDefault();
                      local.setState(() => item);
                      setInput("");
                    }}
                  >
                    {local.getItemName(item)}
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
