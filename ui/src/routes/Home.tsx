import { createBearerSignal, LOGIN_PATH } from "../utils/oidc";

export default function Home() {
  const [bearer, _setBearer] = createBearerSignal();

  if (bearer() === "") {
    window.location.replace(LOGIN_PATH);
    return <p>redirect</p>;
  }

  return <h1>Home</h1>;
}
