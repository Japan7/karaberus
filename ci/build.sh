#!/bin/sh -xe
if [ -n "$CHOST" ]; then
    crossarg="--cross-file ci/$CHOST.txt"
    # HACK: change symlink target so call to qemu-user works
    ln -sf ./libc.so /usr/${CHOST}/usr/lib/ld-musl-${ARCH}.so.1
fi

meson setup /build --buildtype release --strip -Dtest=false -Db_lto=true -Db_lto_mode=thin -Dffmpegaacsucks:default_library=static -Ddakara_check:default_library=static -Db_pie=true -Dffmpeg:programs=disabled -Dffmpeg:tests=disabled -Dffmpeg:encoders=disabled -Dffmpeg:muxers=disabled -Dffmpeg:avfilter=disabled -Dffmpeg:avdevice=disabled -Dffmpeg:postproc=disabled -Dffmpeg:swresample=disabled -Dffmpeg:swscale=disabled -Dffmpeg:decoders=disabled -Dffmpeg:aac_decoder=enabled -Dffmpeg:aac_fixed_decoder=enabled -Dffmpeg:aac_latm_decoder=enabled -Dffmpeg:version3=enabled $crossarg
meson compile -C /build
meson install -C /build --destdir /image --skip-subprojects
