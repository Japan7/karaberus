#!/bin/sh -x
[ -n "$CHOST" ] && crossarg=--cross-file ci/$CHOST.txt
set -e

meson setup /build --buildtype release --strip -Db_lto=true -Db_lto_mode=thin -Dffmpegaacsucks:default_library=static -Ddakara_check:default_library=static -Db_pie=true -Dffmpeg:programs=disabled -Dffmpeg:tests=disabled -Dffmpeg:encoders=disabled -Dffmpeg:muxers=disabled -Dffmpeg:avfilter=disabled -Dffmpeg:avdevice=disabled -Dffmpeg:postproc=disabled -Dffmpeg:swresample=disabled -Dffmpeg:swscale=disabled -Dffmpeg:decoders=disabled -Dffmpeg:aac_decoder=enabled -Dffmpeg:aac_fixed_decoder=enabled -Dffmpeg:aac_latm_decoder=enabled -Dffmpeg:version3=enabled $crossarg
meson compile -C /build
meson install -C /build --destdir /image --skip-subprojects
