import { HiSolidPencil, HiSolidTrash } from "solid-icons/hi";
import { createResource, Index, useContext } from "solid-js";
import AuthorEditor from "../../components/AuthorEditor";
import { Context } from "../../components/Context";
import type { components } from "../../utils/karaberus";
import { karaberus } from "../../utils/karaberus-client";
import { isAdmin } from "../../utils/session";

export default function TagsAuthor() {
  const { getModalRef, setModal, showToast } = useContext(Context)!;

  const [getAuthors, { refetch }] = createResource(async () => {
    const resp = await karaberus.GET("/api/tags/author");
    return resp.data;
  });

  const postAuthor = async (author: components["schemas"]["AuthorInfo"]) => {
    const resp = await karaberus.POST("/api/tags/author", { body: author });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    showToast("Author added!");
    refetch();
  };

  const patchAuthor =
    (id: components["schemas"]["TimingAuthor"]["ID"]) =>
    async (author: components["schemas"]["AuthorInfo"]) => {
      const resp = await karaberus.PATCH("/api/tags/author/{id}", {
        params: { path: { id } },
        body: author,
      });
      if (resp.error) {
        alert(resp.error.detail);
        return;
      }
      getModalRef().close();
      showToast("Author edited!");
      refetch();
    };

  const deleteAuthor = async (
    id: components["schemas"]["TimingAuthor"]["ID"],
  ) => {
    if (!confirm("Confirm deletion?")) {
      return;
    }
    const resp = await karaberus.DELETE("/api/tags/author/{id}", {
      params: { path: { id } },
    });
    if (resp.error) {
      alert(resp.error.detail);
      return;
    }
    showToast("Author deleted!");
    refetch();
  };

  return (
    <>
      <h1 class="header">Author Tags</h1>

      <h2 class="text-2xl font-semibold">Add author</h2>

      <AuthorEditor onSubmit={postAuthor} reset />

      <h2 class="text-2xl font-semibold mt-8">Browse</h2>

      <div class="overflow-auto">
        <table class="table">
          <thead>
            <tr>
              <th></th>
              <th>Name</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <Index each={getAuthors()}>
              {(getAuthor) => (
                <tr class="hover">
                  <th>{getAuthor().ID}</th>
                  <td>{getAuthor().Name}</td>
                  <td class="flex gap-x-1">
                    <button
                      class="btn btn-sm btn-warning"
                      onclick={() => {
                        setModal(
                          <AuthorEditor
                            author={getAuthor()}
                            onSubmit={patchAuthor(getAuthor().ID)}
                          />,
                        );
                        getModalRef().showModal();
                      }}
                    >
                      <HiSolidPencil class="size-4" />
                    </button>
                    <button
                      disabled={!isAdmin()}
                      onclick={() => deleteAuthor(getAuthor().ID)}
                      class="btn btn-sm btn-error"
                    >
                      <HiSolidTrash class="size-4" />
                    </button>
                  </td>
                </tr>
              )}
            </Index>
          </tbody>
        </table>
      </div>
    </>
  );
}
