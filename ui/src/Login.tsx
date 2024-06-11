import { onMount, type Component } from 'solid-js'
import { KARABERUS_API_BASE } from './const'
import { oidc_callback_path } from './oidc'

interface OIDCConfig {
    authorization_endpoint: string
    issuer: string
    jwks_uri: string
    token_endpoint: string
    client_id: string
}

const Login: Component = () => {
    let oidc_discovery_endpoint = KARABERUS_API_BASE + "/oidc_discovery"

    onMount(async () => {
        let oidc_discovery = await fetch(oidc_discovery_endpoint)
        let oidc_config: OIDCConfig = await oidc_discovery.json()

        let auth_endpoint = oidc_config.authorization_endpoint
        let login_url = new URL(auth_endpoint)
        let params = login_url.searchParams

        let callback_uri = `${window.location.protocol}//${window.location.host}${oidc_callback_path}`

        params.append("scope", "openid")
        params.append("response_type", "code")
        params.append("client_id", oidc_config.client_id)
        params.append("redirect_uri", callback_uri)
        // params.append("state", something)

        window.location.replace(login_url)
    })

    return <></>
};

export default Login
