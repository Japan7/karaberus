import { Router } from "@solidjs/router";
import routes from "~solid-pages";
import Layout from "./layout/Layout";

export default function App() {
  return <Router root={Layout}>{routes}</Router>;
}
