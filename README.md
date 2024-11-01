# Karaberus

A karaoke database.

## Getting started

```
meson setup build -Dbuiltin_oidc_env=true -Dbuiltin_s3_env=true
meson compile -C build
```

On Windows this might not build out of the box, for now you can disable the native dependencies which should allow you to build:
```sh
meson setup build --reconfigure -Dno_native_deps=true
```

To run the app you need an oidc server and a s3, for development you can run `meson compile -C build oidc` to have a simple oidc server (`-Dbuiltin_oidc_env=true` has to be set during setup for this to work) and `meson compile -C build s3` to have a S3 server (`-Dbuiltin_s3_env=true` has to be set during setup for this to work).

Then you can start the server with:
```sh
meson compile -C build run
```

This server also serves the web frontend so it should be useable.

To ease development on the web frontend it might be more convenient to run the dev server:
```sh
# from the root of the repo
cd ui
npm run dev
```

# Custom S3 server

If you want to upload files (or run the S3 tests) to your own S3 server you have to set the following environment variables appropriately:

```sh
export KARABERUS_S3_ENDPOINT=localhost:9000
export KARABERUS_S3_KEYID=your_keyid
export KARABERUS_S3_SECRET=your_secret
```

# Custom OIDC server

```sh
export KARABERUS_OIDC_ISSUER=http://localhost:9998/
export KARABERUS_OIDC_KEY_ID=karaberus
export KARABERUS_OIDC_CLIENT_ID=your_client_id
export KARABERUS_OIDC_CLIENT_SECRET=your_client_secret
export KARABERUS_OIDC_GROUPS_CLAIM=groups
export KARABERUS_OIDC_ADMIN_GROUP=admin
export KARABERUS_OIDC_SCOPES="openid profile email groups"
export KARABERUS_OIDC_JWT_SIGN_KEY="" # openssl rand -hex 32
```

## Dakara

You can start a [Dakara server](https://github.com/DakaraProject/dakara-server/) and feed it from your instance:
```sh
export KARABERUS_DAKARA_BASE_URL="http://127.0.0.1:8000"
export KARABERUS_DAKARA_TOKEN="YOUR_DAKARA_TOKEN"
```

## Tests

Tests can easily be run with:
```sh
meson test -C build
```

You can disable the S3 tests with the `s3_tests` configuration option of meson if you have not set up a S3 server:
```sh
meson setup build --reconfigure -Ds3_tests=disabled
```

If you run the tests in a git hook or similar, some tests can be slow to run and can be disabled if needed:
```sh
meson setup build --reconfigure -Dstaticcheck=false -Derrcheck=false
```
(These tests do run in the Github Actions so you should notice it at some point after pushing)

## Development setup

Due to the current Meson setup, this project might not work well out of the box with regular go dev tools.

For the most part `gopls` should work without issues in the `server` module, if you need to work in the `karaberus_tools` (which interfaces with C libraries) you can set `PKG_CONFIG_PATH` to point to `build/meson-uninstalled` in the [settings of gopls](https://github.com/golang/tools/blob/master/gopls/doc/settings.md#env-mapstringstring).

NOTE: You should not set PKG_CONFIG_PATH in your general environment (or at least in the environment in which you run Meson) as it might break meson when it needs to reconfigure itself (when meson.build files are modified most likely), because then it can find its own dependencies with pkg-config and get really confused when it won't find the libraries later in the build.

### Dev Container / Codespaces

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/Japan7/karaberus?quickstart=1)

Once attached to the container, create and edit `.vscode/settings.json` as follows:

1. Copy [`.vscode/settings.example.json`](.vscode/settings.example.json) to `.vscode/settings.json`.

2. For `KARABERUS_LISTEN_BASE_URL`, keep it to `http://localhost:5173` if you are running a Dev Container or a Codespace remotely [attached to your local VS Code](https://docs.github.com/en/codespaces/developing-in-a-codespace/using-github-codespaces-in-visual-studio-code).\
For Codespaces in browser, set it to your editor URL, **without trailing slash**, and **ending with `-5173.app.github.dev` instead of `.github.dev`**.

3. If you want to use `https://auth.japan7.bde.enseeiht.fr` as `KARABERUS_OIDC_ISSUER`, ask for the required `KARABERUS_OIDC_CLIENT_ID` and `KARABERUS_OIDC_CLIENT_SECRET` on Discord. Otherwise, bring your own secrets.

Then use the provided [launch tasks](.vscode/launch.json) to compile and start both the server and the frontend.
