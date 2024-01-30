FROM docker.io/amd64/alpine AS bootstrap
RUN apk add tar gzip curl
RUN curl https://repo.chimera-linux.org/live/20240122/chimera-linux-x86_64-ROOTFS-20240122-bootstrap.tar.gz -o chimera.tar.gz && \
    echo 'b57b4c84ef1b8c5c628e84aa9f1b80863d748440643041414102b9ee19b0a5e4  chimera.tar.gz' | sha256sum -c && \
    mkdir /chimera && tar xf chimera.tar.gz -C /chimera

FROM ghcr.io/japan7/dakara_check:master as dakara_check

FROM scratch AS builder 
COPY --from=bootstrap /chimera /

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
