import { createSignal, onMount } from "solid-js";
import { currentToken, getUserData } from "../utils/oidc";

export default function Home() {
  const [info, setInfo] = createSignal({ foo: "bar" });

  onMount(async () => {
    if (
      currentToken.access_token &&
      currentToken.expires &&
      Date.now() < parseInt(currentToken.expires)
    ) {
      const resp = await getUserData();
      console.log(resp);
      setInfo(resp);
    }
  });

  return (
    <>
      <p>Home</p>
      <pre>{JSON.stringify(info, undefined, 2)}</pre>
    </>
  );
}
