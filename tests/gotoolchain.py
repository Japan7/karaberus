#!/usr/bin/env python
# check that go.mod doesn’t contain a toolchain instruction
import sys
import re

toolchain_prog = re.compile("^toolchain")


def find_toolchain_instruction(file: str):
    with open(file) as f:
        for line in f:
            if toolchain_prog.search(line) is not None:
                print(f"found toolchain instruction in {file}")
                return True

    return False


def main(*files: str):
    if any(map(find_toolchain_instruction, files)):
        return sys.exit(1)


if __name__ == "__main__":
    main(*sys.argv[1:])
