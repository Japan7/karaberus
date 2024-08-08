import createClient from "openapi-fetch";
import type { paths } from "./karaberus";
import { isTauri } from "./tauri";
import { getSessionToken } from "./session";

export const karaberus = createClient<paths>({
  baseUrl: import.meta.env.VITE_KARABERUS_URL || import.meta.env.BASE_URL,
  headers: {
    Authorization: isTauri ? `Bearer ${getSessionToken()}` : undefined,
  },
});

export function apiUrl(path: `api/${string}`) {
  let baseUrl = import.meta.env.VITE_KARABERUS_URL;
  if (baseUrl) {
    baseUrl += "/";
  } else {
    baseUrl = import.meta.env.BASE_URL;
  }
  return baseUrl + path;
}

// Redirect to the home page if the user is not authenticated
karaberus.use({
  onResponse: ({ response }) => {
    if (response.status === 401) {
      // force rerender of authhero
      location.reload();
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
