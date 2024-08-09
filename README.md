# Karaberus

A karaoke database.

## Getting started

```
meson setup build
meson compile -C build
```

On Windows this might not build out of the box, for now you can disable the native dependencies which should allow you to build:
```sh
meson setup build --reconfigure -Dno_native_deps=true
```

To run the app you need an oidc server, for development [zitadel-karaberus](https://github.com/odrling/zitadel-karaberus) can be used.

If you want to upload files (or run the S3 tests) you will also need an S3 server, probably the easiest way to have a local server is to use the Minio Docker image:

```sh
podman run -d -e MINIO_DEFAULT_BUCKETS=karaberus -p 9000:9000 bitnami/minio:latest
```

Then for the environment variables you have to set:
```sh
export KARABERUS_S3_ENDPOINT=localhost:9000
export KARABERUS_S3_KEYID=minio
export KARABERUS_S3_SECRET=miniosecret
```

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

# Development setup

Due to the current Meson setup, this project might not work well out of the box with regular go dev tools.

For the most part `gopls` should work without issues in the `server` module, if you need to work in the `karaberus_tools` (which interfaces with C libraries) you can set `PKG_CONFIG_PATH` to point to `build/meson-uninstalled` in the [settings of gopls](https://github.com/golang/tools/blob/master/gopls/doc/settings.md#env-mapstringstring).

NOTE: You should not set PKG_CONFIG_PATH in your general environment as it might break meson when it needs to reconfigure itself (when meson.build files are modified most likely), because then it can find its own dependencies with pkg-config and get really confused when it won't find the libraries later in the build.

