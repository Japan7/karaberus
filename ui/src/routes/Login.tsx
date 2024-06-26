import { onMount } from "solid-js";
import { karaberus } from "../utils/karaberus-client";
import { OIDC_CALLBACK_PATH } from "../utils/oidc";

export default function Login() {
  onMount(async () => {
    const { data: oidc_config, error } = await karaberus.GET(
      "/api/oidc_discovery",
    );

    if (error) {
      alert(error.detail);
      throw error;
    }

    const auth_endpoint = oidc_config.authorization_endpoint;
    const login_url = new URL(auth_endpoint);
    const params = login_url.searchParams;

    const callback_uri = `${window.location.protocol}//${window.location.host}${OIDC_CALLBACK_PATH}`;

    params.append("scope", "openid");
    params.append("response_type", "code");
    params.append("client_id", oidc_config.client_id);
    params.append("redirect_uri", callback_uri);
    // params.append("state", something)

    window.location.replace(login_url);
  });

  return <></>;
}
