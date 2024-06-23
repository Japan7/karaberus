#!/usr/bin/env python3
import sys
import pathlib
import subprocess
import os


def run(*args: str):
    return subprocess.run(args, check=True)


def main(*args: str):
    PKG_CONFIG = os.environ.get("DEVENV_PKG_CONFIG", "pkg-config")
    MESON_BUILD_DIR = os.environ["MESON_BUILD_DIR"]
    builddir_path = pathlib.Path(MESON_BUILD_DIR)

    if "karaberus_tools" in args and builddir_path.exists():
        os.environ["PKG_CONFIG_PATH"] = str(builddir_path / "meson-uninstalled")

    return run(PKG_CONFIG, *args)


if __name__ == "__main__":
    main(*sys.argv[1:])
