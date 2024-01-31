FROM ghcr.io/japan7/dakara_check:master as dakara_check

FROM ghcr.io/odrling/chimera-images:main AS builder 

RUN apk add chimera-repo-contrib

ARG ARCH
ARG GOARCH
ARG SYSROOT

COPY --from=dakara_check /usr/local $SYSROOT/usr/

COPY . /karaberus

RUN cd /karaberus && ARCH=$ARCH GOARCH=$GOARCH SYSROOT=$SYSROOT ci/build.sh

FROM alpine

COPY --from=builder /karaberus/build/karaberus /
ENV KARABERUS_LISTEN_ADDR=":8888"
EXPOSE 8888
ENTRYPOINT ["/karaberus"]
