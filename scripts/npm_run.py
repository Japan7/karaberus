import subprocess
import os
import sys

def run(*args: str):
    return subprocess.run(args, check=True)

def main():
    npm = os.environ["NPM"]
    ui_source = os.environ["SOURCE"]

    npm_cmd = npm, "-C", ui_source
    run(*npm_cmd, "ci")

    npm_subcmd = sys.argv[1]
    if npm_subcmd == "build":
        outdir = sys.argv[2]
        run(*npm_cmd, "run", npm_subcmd, "--", "--outDir", outdir)
    else:
        run(*npm_cmd, "run", npm_subcmd)


if __name__ == "__main__":
    main()
