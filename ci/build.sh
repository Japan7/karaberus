#!/bin/sh -xe
export GOCACHE=/karaberus/go_cache
export GOTOOLCHAIN=local

if [ -n "${TARGET}" ]; then
    crossarg="--cross-file /karaberus/ci/${TARGET}.ini"
    SYSROOT="/usr/${TARGET}"
fi

IMAGE=/karaberus/image

mkdir -p ${IMAGE}/etc
cp -r /etc/ssl ${IMAGE}/etc/

meson setup /build /karaberus --buildtype release --strip --libdir lib --prefix ${IMAGE} -Db_lto=true -Db_lto_mode=thin -Db_pie=true -Dc_args=-fhardened $crossarg
meson test -C /build --verbose
meson install -C /build --tags runtime

mkdir ${IMAGE}/tmp
chmod 1777 ${IMAGE}/tmp

SYSROOT=${SYSROOT} /karaberus/ci/copyimglibs.py ${IMAGE} ${IMAGE}/bin/karaberus
if [ -d "${SYSROOT}/lib64" ]; then
    mkdir -p ${IMAGE}/lib64
    ln ${IMAGE}/lib/ld-* ${IMAGE}/lib64/
fi
