import { onMount } from "solid-js";
import { redirectToAuthorize } from "../utils/oidc";

export default function OIDCLogin() {
  onMount(redirectToAuthorize);

  return null;
}
