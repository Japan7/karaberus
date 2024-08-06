import { Router } from "@solidjs/router";
import { createEffect } from "solid-js";
import { themeChange } from "theme-change";
import routes from "~solid-pages";
import Layout from "./layout/Layout";

export default function App() {
  createEffect(() => themeChange());

  return <Router root={Layout}>{routes}</Router>;
}
