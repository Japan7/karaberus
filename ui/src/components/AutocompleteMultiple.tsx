import { HiSolidXMark } from "solid-icons/hi";
import { createSignal, For, type Accessor, type Setter } from "solid-js";

export default function AutocompleteMultiple<T>({
  items,
  getItemName,
  getState,
  setState,
  required,
  allowDuplicates,
}: {
  items: T[];
  getItemName: (item: T) => string;
  getState: Accessor<T[]>;
  setState: Setter<T[]>;
  required?: boolean;
  allowDuplicates?: boolean;
}) {
  let inputEl!: HTMLInputElement;
  const [getInput, setInput] = createSignal("");

  const filteredItems = () =>
    items.filter(
      (item) =>
        (allowDuplicates || !getState().includes(item)) &&
        getItemName(item).toLowerCase().includes(getInput().toLowerCase()),
    );

  return (
    <div class="dropdown w-full">
      <div class="textarea textarea-bordered w-full">
        <div class="flex flex-wrap gap-1">
          <For each={getState()}>
            {(item, index) => (
              <div
                onclick={() =>
                  setState((state) => state.filter((_, i) => i !== index()))
                }
                class="btn btn-xs"
              >
                <span>{getItemName(item)}</span>
                <HiSolidXMark class="size-4" />
              </div>
            )}
          </For>
        </div>
        <input
          type="text"
          required={required && getState().length === 0}
          value={getInput()}
          oninput={(e) => setInput(e.currentTarget.value)}
          ref={inputEl}
          class="outline-none w-full"
        />
      </div>
      <ul class="dropdown-content menu bg-base-100 rounded-box z-[1] w-full p-2 shadow">
        <For each={filteredItems()}>
          {(item) => (
            <li>
              <a
                href="#"
                onclick={(e) => {
                  e.preventDefault();
                  setState((state) => [...state, item]);
                  setInput("");
                  inputEl.focus();
                }}
              >
                {getItemName(item)}
              </a>
            </li>
          )}
        </For>
      </ul>
    </div>
  );
}
