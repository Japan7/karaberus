#!/bin/sh
meson setup -Db_lto=true -Db_lto_mode=thin -Db_pie=true -Dc_args=-fhardened /build /karaberus
meson test -C /build --verbose
