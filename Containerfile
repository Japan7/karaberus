FROM ghcr.io/japan7/dakara_check:master

RUN apk add go pkgconf clang llvm lld

COPY . /karaberus

RUN cd /karaberus && CGO_ENABLED=1 CC=clang go build -o build/ .
