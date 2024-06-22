import createClient from "openapi-fetch";
import type { paths } from "./karaberus";

export const karaberus = createClient<paths>({
  baseUrl: import.meta.env.BASE_URL,
});
