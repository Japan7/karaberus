import { onMount } from "solid-js";
import { currentToken, getToken } from "../utils/oidc";
import routes from "../utils/routes";

export default function OIDCCallback() {
  onMount(async () => {
    const args = new URLSearchParams(window.location.search);
    const code = args.get("code");

    // If we find a code, we're in a callback, do a token exchange
    if (code) {
      const token = await getToken(code);
      currentToken.save(token);
    }

    window.location.href = routes.HOME;
  });

  return null;
}
