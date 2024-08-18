import { debounce } from "@solid-primitives/scheduled";
import Fuse from "fuse.js";
import { HiSolidXMark } from "solid-icons/hi";
import {
  createEffect,
  createSignal,
  For,
  splitProps,
  type Accessor,
  type JSX,
  type Setter,
} from "solid-js";

export default function AutocompleteMultiple<T>(
  props: {
    items: T[];
    getItemName: (item: T) => string;
    getState: Accessor<T[]>;
    setState: Setter<T[]>;
    allowDuplicates?: boolean;
  } & JSX.InputHTMLAttributes<HTMLInputElement>,
) {
  const [local, inputProps] = splitProps(props, [
    "items",
    "getItemName",
    "getState",
    "setState",
    "allowDuplicates",
  ]);

  const [getInput, setInput] = createSignal("");
  const [getFilteredItems, setFilteredItems] = createSignal(local.items);

  let inputRef!: HTMLInputElement;
  let dropdownRef!: HTMLDivElement;

  const filter = (query: string) => {
    const filtered = local.items.filter(
      (item) => local.allowDuplicates || !local.getState().includes(item),
    );
    if (!query) {
      return filtered;
    }
    const fuse = new Fuse(filtered, {
      keys: [
        {
          name: "name",
          getFn: (item) => local.getItemName(item),
        },
      ],
    });
    return fuse.search(query).map((result) => result.item);
  };
  const debouncedFilter = debounce((query: string) => {
    setFilteredItems(filter(query));
  }, 250);

  createEffect(() => {
    debouncedFilter(getInput());
  });

  const handleKeyDownInput: JSX.EventHandler<
    HTMLInputElement,
    KeyboardEvent
  > = (e) => {
    if (e.key === "ArrowUp" || e.key === "ArrowDown") {
      e.preventDefault();
      const items = dropdownRef.querySelectorAll("a");
      items[e.key === "ArrowUp" ? items.length - 1 : 0].focus();
    }
    if (
      e.key === "Backspace" &&
      e.currentTarget.selectionStart === 0 &&
      e.currentTarget.selectionEnd === 0
    ) {
      e.preventDefault();
      local.setState((state) => state.slice(0, -1));
    }
  };

  const handleKeyDownDropdown: JSX.EventHandler<
    HTMLDivElement,
    KeyboardEvent
  > = (e) => {
    if (e.key === "ArrowUp" || e.key === "ArrowDown") {
      e.preventDefault();
      const items = dropdownRef.querySelectorAll("a");
      const index = Array.from(items).indexOf(
        document.activeElement as HTMLAnchorElement,
      );
      if (
        (index === 0 && e.key === "ArrowUp") ||
        (index === items.length - 1 && e.key === "ArrowDown")
      ) {
        inputRef.focus();
      } else {
        items[index + (e.key === "ArrowUp" ? -1 : 1)].focus();
      }
    }
  };

  return (
    <div class="dropdown w-full">
      <div class="textarea textarea-bordered w-full">
        <div class="flex flex-wrap gap-1">
          <For each={local.getState()}>
            {(item, index) => (
              <div
                onclick={() =>
                  local.setState((state) =>
                    state.filter((_, i) => i !== index()),
                  )
                }
                class="btn btn-xs"
              >
                <span>{local.getItemName(item)}</span>
                <HiSolidXMark class="size-4" />
              </div>
            )}
          </For>
        </div>
        <input
          type="text"
          required={inputProps.required && local.getState().length === 0}
          value={getInput()}
          oninput={(e) => setInput(e.currentTarget.value)}
          onkeydown={handleKeyDownInput}
          ref={inputRef}
          class="bg-transparent outline-none w-full"
          tabindex="0"
          {...inputProps}
        />
      </div>
      <div
        onkeydown={handleKeyDownDropdown}
        ref={dropdownRef}
        class="dropdown-content menu bg-base-100 rounded-box z-[1] w-full p-2 shadow max-h-48 overflow-y-auto"
        tabindex="0"
      >
        <ul>
          <For each={getFilteredItems()}>
            {(item) => (
              <li>
                <a
                  href="#"
                  onclick={(e) => {
                    e.preventDefault();
                    local.setState((state) => [...state, item]);
                    setInput("");
                    inputRef.focus();
                  }}
                >
                  {local.getItemName(item)}
                </a>
              </li>
            )}
          </For>
        </ul>
      </div>
    </div>
  );
}
