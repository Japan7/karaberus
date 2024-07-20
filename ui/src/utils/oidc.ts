// Adapted from
// https://github.com/spotify/web-api-examples/blob/7c4872d343a6f29838c437cf163012947b4bffb9/authorization/authorization_code_pkce/public/app.js
import { karaberus } from "./karaberus-client";
import routes from "./routes";

export interface TokenResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  refresh_token: string;
  scope: string;
}

// Data structure that manages the current active token, caching it in localStorage
export const currentToken = {
  get access_token() {
    return localStorage.getItem("access_token");
  },
  get refresh_token() {
    return localStorage.getItem("refresh_token");
  },
  get expires_in() {
    const expiresIn = localStorage.getItem("expires_in");
    return expiresIn ? parseInt(expiresIn) : null;
  },
  get expires_at() {
    const expiresAt = localStorage.getItem("expires_at");
    return expiresAt ? parseInt(expiresAt) : null;
  },
  get isValid() {
    return this.access_token && this.expires_at && Date.now() < this.expires_at;
  },

  save: function (resp: TokenResponse) {
    const { access_token, refresh_token, expires_in } = resp;
    localStorage.setItem("access_token", access_token);
    localStorage.setItem("refresh_token", refresh_token);
    localStorage.setItem("expires_in", expires_in.toString());

    const expiresAt = Date.now() + expires_in * 1000;
    localStorage.setItem("expires_at", expiresAt.toString());
  },
};

export async function redirectToAuthorize() {
  const possible =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~";
  const randomValues = crypto.getRandomValues(new Uint8Array(64));
  const randomString = randomValues.reduce(
    (acc, x) => acc + possible[x % possible.length],
    "",
  );

  const code_verifier = randomString;
  const data = new TextEncoder().encode(code_verifier);
  const hashed = await crypto.subtle.digest("SHA-256", data);

  const code_challenge_base64 = btoa(
    String.fromCharCode(...new Uint8Array(hashed)),
  )
    .replace(/=/g, "")
    .replace(/\+/g, "-")
    .replace(/\//g, "_");

  window.localStorage.setItem("code_verifier", code_verifier);

  const { data: oidc_config, error } = await karaberus.GET(
    "/api/oidc_discovery",
  );
  if (error) {
    throw new Error(error.detail);
  }

  const authUrl = new URL(oidc_config.authorization_endpoint);
  const params = {
    response_type: "code",
    client_id: oidc_config.client_id,
    scope: "openid profile email",
    code_challenge_method: "S256",
    code_challenge: code_challenge_base64,
    redirect_uri: `${window.location.protocol}//${window.location.host}${routes.OIDC_CALLBACK}`,
  };

  authUrl.search = new URLSearchParams(params).toString();
  window.location.href = authUrl.toString(); // Redirect the user to the authorization server for login
}

export async function getToken(code: string) {
  const code_verifier = localStorage.getItem("code_verifier");

  const { data: oidc_config, error } = await karaberus.GET(
    "/api/oidc_discovery",
  );
  if (error) {
    throw new Error(error.detail);
  }

  const resp = await fetch(oidc_config.token_endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      client_id: oidc_config.client_id,
      grant_type: "authorization_code",
      code: code,
      redirect_uri: `${window.location.protocol}//${window.location.host}${routes.OIDC_CALLBACK}`,
      code_verifier: code_verifier!,
    }),
  });

  const json = await resp.json();
  return json as TokenResponse;
}

export async function refreshToken() {
  const { data: oidc_config, error } = await karaberus.GET(
    "/api/oidc_discovery",
  );
  if (error) {
    throw new Error(error.detail);
  }

  const resp = await fetch(oidc_config.token_endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      client_id: oidc_config.client_id,
      grant_type: "refresh_token",
      refresh_token: currentToken.refresh_token!,
    }),
  });

  const json = await resp.json();
  return json as TokenResponse;
}
