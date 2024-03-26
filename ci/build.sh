#!/bin/sh -x
if [ "$ARCH" != x86_64 ]; then
    CC=$ARCH-chimera-linux-musl-clang
    mv "$SYSROOT/usr/local/lib/"* "$SYSROOT/usr/lib/"
else
    CC=clang
fi

export CC

set -e
apk add chimera-repo-contrib
apk add go

export PKG_CONFIG_PATH="$SYSROOT/usr/lib/pkgconfig"
export CGO_ENABLED=1
export GOARCH="$GOARCH"
export GOOS=linux
export CGO_CFLAGS="-fPIE -O3 -Wall -Wextra"

go build -buildmode=pie -trimpath -ldflags '-linkmode=external -s' -o build/ .
