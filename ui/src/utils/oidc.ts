import { makePersisted } from "@solid-primitives/storage";
import { Signal, createSignal } from "solid-js";

export function createBearerSignal(): Signal<string> {
  return makePersisted(createSignal(""), {
    name: "bearer",
    storage: sessionStorage,
  });
}

export const LOGIN_PATH = "/login";
export const OIDC_CALLBACK_PATH = "/oidc/callback";
