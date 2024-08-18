#!/bin/sh -e
if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=$1

git switch master
yq ".version = \"${VERSION}\"" -i ui/src-tauri/tauri.conf.json
git add ui/src-tauri/tauri.conf.json
git commit -m "desktop: release ${VERSION}"
git tag app-v${VERSION}
git push --atomic origin master app-v${VERSION}
git switch -
