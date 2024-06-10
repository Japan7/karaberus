import type { Component } from 'solid-js';
import { createBearerSignal } from './oidc';

const OIDCCallback: Component = () => {
    const [bearer, setBearer] = createBearerSignal()

    let uri = new URL(window.location.href)
    let code = uri.searchParams.get("code")

    if (code !== null) {
        setBearer(code!)
    }

    window.location.replace("/")
    return <></>;
};

export default OIDCCallback;
