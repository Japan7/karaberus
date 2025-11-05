import {
  createResource,
  createSignal,
  Show,
  useContext,
  type JSX,
} from "solid-js";
import type { components } from "../utils/karaberus";
import { karaberus } from "../utils/karaberus-client";
import AuthorEditor from "./AuthorEditor";
import Autocomplete from "./Autocomplete";
import { Context } from "./Context";

export default function ProfileEditor(props: {
  user: components["schemas"]["User"];
  onSubmit: (data: { author?: components["schemas"]["TimingAuthor"] }) => void;
  reset?: boolean;
}) {
  const { getModalRef, setModal, showToast } = useContext(Context)!;

  const [getAllAuthors, { refetch: refetchAuthors }] = createResource(
    async () => {
      const resp = await karaberus.GET("/api/tags/author");
      return resp.data;
    },
  );

  const [getAuthor, setAuthor] = createSignal<
    components["schemas"]["TimingAuthor"] | undefined
  >(props.user.timing_profile);

  const openAddAuthorModal: JSX.EventHandler<HTMLElement, MouseEvent> = (e) => {
    e.preventDefault();
    setModal(<AuthorEditor onSubmit={postAuthor} />);
    getModalRef().showModal();
  };

  const postAuthor = async (author: components["schemas"]["AuthorInfo"]) => {
    const resp = await karaberus.POST("/api/tags/author", { body: author });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    showToast("Author added!");
    refetchAuthors();
    getModalRef().close();
  };

  const onsubmit: JSX.EventHandler<HTMLElement, SubmitEvent> = (e) => {
    e.preventDefault();
    props.onSubmit({
      author: getAuthor(),
    });
    if (props.reset) {
      (e.target as HTMLFormElement).reset();
    }
  };

  return (
    <form onsubmit={onsubmit} class="flex flex-col gap-y-2 w-full">
      <label>
        <div class="label">
          <span>Timing Author</span>
          <span class="text-sm opacity-70">
            <a class="link" onclick={openAddAuthorModal}>
              Can't find it?
            </a>
          </span>
        </div>
        <Show
          when={getAllAuthors()}
          fallback={<span class="loading loading-spinner loading-lg" />}
        >
          {(getAllAuthors) => (
            <Autocomplete
              items={getAllAuthors()}
              getItemName={(author) => author.Name}
              getState={getAuthor}
              setState={setAuthor}
              placeholder="bebou69"
            />
          )}
        </Show>
      </label>

      <input type="submit" class="btn btn-primary" />
    </form>
  );
}
