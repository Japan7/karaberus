import createClient from "openapi-fetch";
import type { paths } from "./karaberus";

export const karaberus = createClient<paths>({
  baseUrl: import.meta.env.BASE_URL,
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
