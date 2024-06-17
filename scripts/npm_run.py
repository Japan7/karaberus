import os
import pathlib
import subprocess
import sys

def run(*args: str):
    return subprocess.run(args, check=True)

def main():
    npm = os.environ["NPM"]
    ui_source = os.environ["SOURCE"]

    npm_cmd = npm, "-C", ui_source

    npm_subcmd = sys.argv[1]
    if npm_subcmd == "build":
        outdir = sys.argv[2]
        run(*npm_cmd, "run", npm_subcmd, "--", "--outDir", outdir)
    elif npm_subcmd == "ci":
        run(*npm_cmd, "ci")
        outpath = sys.argv[2]
        with open(outpath, "w"):
            pass
    else:
        run(*npm_cmd, "run", npm_subcmd)


if __name__ == "__main__":
    main()
