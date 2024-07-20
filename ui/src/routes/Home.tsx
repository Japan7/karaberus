import { createSignal } from "solid-js";

export default function Home() {
  const [info, setInfo] = createSignal();

  return (
    <>
      <p>Home</p>
    </>
  );
}
