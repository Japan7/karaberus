import type { Component } from 'solid-js';
import { createBearerSignal, login_path } from './oidc';

const Home: Component = () => {
  const [bearer, _] = createBearerSignal()

  if (bearer() === "") {
    window.location.replace(login_path)
    return <p>redirect</p>
  }

  return <h1>Home</h1>;
};

export default Home;
