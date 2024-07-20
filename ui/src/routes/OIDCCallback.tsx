import { onMount } from "solid-js";
import { currentToken, getToken } from "../utils/oidc";

export default function OIDCCallback() {
  onMount(async () => {
    const args = new URLSearchParams(window.location.search);
    const code = args.get("code");

    // If we find a code, we're in a callback, do a token exchange
    if (code) {
      const token = await getToken(code);
      currentToken.save(token);

      // Remove code from URL so we can refresh correctly.
      const url = new URL(window.location.href);
      url.searchParams.delete("code");

      const updatedUrl = url.search ? url.href : url.href.replace("?", "");
      window.history.replaceState({}, document.title, updatedUrl);
    }
  });
  return null;
}
