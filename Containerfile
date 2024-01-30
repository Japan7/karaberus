FROM ghcr.io/japan7/dakara_check:master

RUN apk add go pkgconf

COPY . /karaberus

RUN cd /karaberus && go build -o build/ .
