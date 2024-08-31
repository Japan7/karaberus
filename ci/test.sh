#!/bin/sh
meson setup /build /karaberus
meson test -C /build --verbose
