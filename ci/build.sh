#!/bin/sh -xe
if [ -n "${TARGET}" ]; then
    crossarg="--cross-file ci/${TARGET}.ini"
fi

meson setup build --buildtype release --strip -Dtest=false -Ds3_tests=disabled -Db_lto=true -Db_lto_mode=thin -Db_pie=true -Dffmpeg:programs=disabled -Dffmpeg:tests=disabled -Dffmpeg:encoders=disabled -Dffmpeg:muxers=disabled -Dffmpeg:avfilter=disabled -Dffmpeg:avdevice=disabled -Dffmpeg:postproc=disabled -Dffmpeg:swresample=disabled -Dffmpeg:swscale=disabled -Dffmpeg:decoders=disabled -Dffmpeg:aac_decoder=enabled -Dffmpeg:aac_fixed_decoder=enabled -Dffmpeg:aac_latm_decoder=enabled -Dffmpeg:version3=enabled $crossarg
meson install -C build --destdir image --tags runtime
ccache --show-stats
