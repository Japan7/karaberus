#!/bin/sh -x
pkgs="musl-devel meson clang clang-rt-devel llvm-devel-static git lld go pkgconf"
if [ "$ARCH" != x86_64 ]; then
    pkgs="$pkgs base-cross-$ARCH"
    export CC=$ARCH-chimera-linux-musl-clang
    mv $SYSROOT/usr/local/lib/* $SYSROOT/usr/lib/
else
    export CC=clang
fi

set -e
apk add $pkgs

export PKG_CONFIG_PATH=$SYSROOT/usr/lib/pkgconfig
export CGO_ENABLED=1
export GOARCH=$GOARCH
export GOOS=linux
#export CGO_CFLAGS="-print-search-dirs"

go build -ldflags '-linkmode=external -s' -o build/ .
