#!/bin/sh -xe
export GOTOOLCHAIN=local

if [ -n "${TARGET}" ]; then
    crossarg="--cross-file /karaberus/ci/${TARGET}.ini"
    SYSROOT="/usr/${TARGET}"
fi

IMAGE=/karaberus/image
BUILDDIR=/karaberus/build

mkdir -p ${IMAGE}/etc
cp -r /etc/ssl ${IMAGE}/etc/

meson setup --reconfigure "${BUILDDIR}" /karaberus --buildtype release -Dbuiltin_oidc_env=false -Dbuiltin_s3_env=false --strip -Dnetwork_tests=disabled --libdir lib --prefix ${IMAGE} -Db_lto=true -Db_lto_mode=thin -Db_pie=true -Dc_args=-fhardened $crossarg
[ "$1" = "--tests" ] && meson test -C "${BUILDDIR}" --verbose
meson install -C "${BUILDDIR}" --tags runtime

mkdir ${IMAGE}/tmp
chmod 1777 ${IMAGE}/tmp

SYSROOT=${SYSROOT} /karaberus/ci/copyimglibs.py ${IMAGE} ${IMAGE}/bin/karaberus
if [ -d "${SYSROOT}/lib64" ]; then
    mkdir -p ${IMAGE}/lib64
    ln ${IMAGE}/lib/ld-* ${IMAGE}/lib64/
fi

golangci-lint cache status
