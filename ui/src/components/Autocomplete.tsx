import { debounce } from "@solid-primitives/scheduled";
import Fuse from "fuse.js";
import { HiSolidXMark } from "solid-icons/hi";
import {
  createEffect,
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
  const [getFilteredItems, setFilteredItems] = createSignal(local.items);

  let inputRef!: HTMLInputElement;
  let dropdownRef!: HTMLDivElement;

  const filter = (query: string) => {
    if (!query) {
      return local.items;
    }
    const fuse = new Fuse(local.items, {
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
        <Show
          when={local.getState()}
          fallback={
            <input
              type="text"
              value={getInput()}
              oninput={(e) => setInput(e.currentTarget.value)}
              onkeydown={handleKeyDownInput}
              ref={inputRef}
              class="bg-transparent outline-none w-full"
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
        <div
          onkeydown={handleKeyDownDropdown}
          ref={dropdownRef}
          class="dropdown-content menu bg-base-100 rounded-box z-[1] w-full p-2 shadow max-h-48 overflow-y-auto"
        >
          <ul>
            <For each={getFilteredItems()}>
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
