import { createContext, createSignal, Show, type JSX } from "solid-js";

export const Context = createContext<{
  getModalRef: () => HTMLDialogElement;
  setModal: (form: JSX.Element) => void;
  showToast: (msg: string) => void;
}>();

export function Provider(props: { children?: JSX.Element }) {
  let modalRef!: HTMLDialogElement;
  const [getModal, setModal] = createSignal<JSX.Element>();

  const [getToast, setToast] = createSignal<string>();

  const showToast = (msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(), 3000);
  };

  return (
    <Context.Provider
      value={{
        getModalRef: () => modalRef,
        setModal,
        showToast,
      }}
    >
      {props.children}

      <dialog ref={modalRef} class="modal modal-bottom sm:modal-middle">
        <div class="modal-box flex justify-center">{getModal()}</div>
        <form method="dialog" class="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>

      <Show when={getToast()}>
        {(getToast) => (
          <div class="toast">
            <div class="alert alert-success">
              <span>{getToast()}</span>
            </div>
          </div>
        )}
      </Show>
    </Context.Provider>
  );
}
