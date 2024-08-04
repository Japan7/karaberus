import { getSessionInfos } from "../utils/session";

export default function Profile() {
  const infos = getSessionInfos();

  return (
    <>
      <h1 class="header">Profile</h1>

      <pre>{JSON.stringify(infos, undefined, 2)}</pre>
    </>
  );
}
