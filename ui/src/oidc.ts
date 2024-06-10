import { makePersisted } from "@solid-primitives/storage";
import { Signal, createSignal } from "solid-js";

function createBearerSignal(): Signal<string> {
    return makePersisted(createSignal(""), { name: "bearer", storage: sessionStorage});
}

const login_path = "/login"
const oidc_callback_path = "/oidc/callback"

export {
    createBearerSignal,
    login_path,
    oidc_callback_path
}
