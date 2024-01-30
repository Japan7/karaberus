FROM ghcr.io/japan7/dakara_check:master

RUN apk add go pkgconf

COPY . /karaberus

RUN cd /karaberus && CGO_ENABLED=1 go build -o build/ .
