#!/usr/bin/env python3
import pathlib
import sys


def main(src: pathlib.Path, dest: pathlib.Path):
    dest.unlink(missing_ok=True)
    dest.symlink_to(src.absolute())


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(f"USAGE: {sys.argv[0]} <src> <dest>")
        sys.exit(1)

    main(pathlib.Path(sys.argv[1]), pathlib.Path(sys.argv[2]))
