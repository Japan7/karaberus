import { onMount } from "solid-js";
import { OIDC_CALLBACK_PATH } from "../utils/oidc";

interface OIDCConfig {
  authorization_endpoint: string;
  issuer: string;
  jwks_uri: string;
  token_endpoint: string;
  client_id: string;
}

export default function Login() {
  onMount(async () => {
    const oidc_discovery = await fetch("/api/oidc_discovery");
    const oidc_config: OIDCConfig = await oidc_discovery.json();

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
