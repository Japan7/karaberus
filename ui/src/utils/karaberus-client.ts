import createClient from "openapi-fetch";
import type { paths } from "./karaberus";
import routes from "./routes";

export const karaberus = createClient<paths>({
  baseUrl: import.meta.env.BASE_URL,
});

// Redirect to the home page if the user is not authenticated
karaberus.use({
  onResponse: ({ response }) => {
    if (response.status === 401) {
      location.href = routes.HOME;
    }
  },
});

export function fileForm(file: File): Record<string, unknown> {
  return {
    body: file,
    bodySerializer: (body: File) => {
      const formData = new FormData();
      formData.append("file", body);
      return formData;
    },
  };
}
