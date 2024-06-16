import { createBearerSignal } from "../utils/oidc";

export default function OIDCCallback() {
  const [_bearer, setBearer] = createBearerSignal();

  const uri = new URL(window.location.href);
  const code = uri.searchParams.get("code");

  if (code !== null) {
    setBearer(code!);
  }

  window.location.replace("/");
  return <></>;
}
