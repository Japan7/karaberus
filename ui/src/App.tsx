import { Router } from "@solidjs/router";
import { isTauri } from "@tauri-apps/api/core";
import { createEffect } from "solid-js";
import { themeChange } from "theme-change";
import routes from "~solid-pages";
import Layout from "./layout/Layout";
import { registerGlobalListeners } from "./utils/tauri";

export default function App() {
  if (isTauri()) {
    registerGlobalListeners();
  }

  createEffect(() => themeChange());

  return <Router root={Layout}>{routes}</Router>;
}
