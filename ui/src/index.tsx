/* @refresh reload */
import { render } from "solid-js/web";

import { isTauri } from "@tauri-apps/api/core";
import App from "./App";
import "./index.css";

if (isTauri()) {
  import("./index.tauri.css");
}

const root = document.getElementById("root");

render(() => <App />, root!);
