#!/usr/bin/env python3
import os
import pathlib
import re
import subprocess
import sys
import shutil
from dataclasses import dataclass


@dataclass
class ObjDumpResult:
    shared_libs: list[str]
    runpath: list[pathlib.Path]


objdump_parser = re.compile(r"^\s*(NEEDED|RUNPATH)\s*(\S*)\s*$")


def objdump(file: pathlib.Path) -> ObjDumpResult:
    objdump_bin = os.environ.get("OBJDUMP", "objdump")
    cmd = [objdump_bin, "-p", str(file)]
    proc = subprocess.run(cmd, stdout=subprocess.PIPE)

    runpath: list[pathlib.Path] = []
    shared_libs: list[str] = []
    for line in proc.stdout.splitlines():
        lineparse = objdump_parser.search(line.decode())
        if lineparse is not None:
            typekey, value = lineparse.groups()
            if typekey == "NEEDED":
                shared_libs.append(value)
            elif typekey == "RUNPATH":
                runpath = list(
                    map(
                        pathlib.Path,
                        value.replace("$ORIGIN", str(file.parent)).split(":"),
                    )
                )

    return ObjDumpResult(shared_libs, runpath)


def findlib(shared_name: str, library_path: list[pathlib.Path]) -> pathlib.Path:
    for libdir in library_path:
        possible_shared_obj = libdir / shared_name
        if possible_shared_obj.is_file():
            return possible_shared_obj

    raise RuntimeError(f"could not find {shared_name} in {library_path}")


def get_sysroot():
    return pathlib.Path(os.environ.get("SYSROOT", "/"))


def parse_ld_so_conf(file: pathlib.Path) -> list[pathlib.Path]:
    if not file.exists():
        return []

    # read ld.so.conf, ignoring includes
    libpath: list[pathlib.Path] = []
    with file.open() as conf:
        for line in conf:
            if line.startswith("/"):
                libdir = pathlib.Path(line.strip())
                # append libdir to prefix
                prefixed_libdir = get_sysroot().joinpath(*libdir.parts[1:])
                libpath.append(prefixed_libdir)
    return libpath


def get_library_path(runpath: list[pathlib.Path], image_dir: pathlib.Path):
    libpath = runpath.copy()
    if library_path_env := os.environ.get("LIBRARY_PATH"):
        libpath.extend(map(pathlib.Path, library_path_env.split(":")))

    ld_so_conf = get_sysroot() / "etc" / "ld.so.conf"
    libpath.extend(parse_ld_so_conf(ld_so_conf))

    # read files in /etc/ld.so.conf.d (since we don't parse includes but it is generally there)
    ld_so_conf_d = get_sysroot() / "etc" / "ld.so.conf.d"
    if ld_so_conf_d.is_dir():
        for ld_so_conf in ld_so_conf_d.iterdir():
            libpath.extend(parse_ld_so_conf(ld_so_conf))

    if get_sysroot() / "lib" not in libpath:
        libpath.append(get_sysroot() / "lib")

    libpath.append(image_dir / "lib")

    return libpath


def find_related_files(file: pathlib.Path, image_dir: pathlib.Path, known_libs: set[str]):
    res = objdump(file)
    libpath = get_library_path(res.runpath, image_dir)
    libs = {findlib(shared_name, libpath) for shared_name in res.shared_libs
            if shared_name not in known_libs}
    for lib in libs.copy():
        libs.update(find_related_files(lib, image_dir, {lib.name for lib in libs}))

    return libs


def main(dest_dir: pathlib.Path, *files: pathlib.Path):
    dest_dir.mkdir(parents=True, exist_ok=True)

    libs_to_copy: set[pathlib.Path] = set()

    for file in files:
        libs_to_copy.update(find_related_files(file, dest_dir, set()))

    for file in libs_to_copy:
        dest_file = dest_dir / "lib" / file.name
        if dest_file.exists():
            print(f"{dest_file} already exists")
            continue

        shutil.copy(file, dest_file)
        print(f"{file} → {dest_file}")


if __name__ == "__main__":
    args = map(pathlib.Path, sys.argv[1:])
    main(*args)
