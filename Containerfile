FROM ghcr.io/japan7/dakara_check:master@sha256:6febecdd7d9d1ffc68a02932b7f879e2ede980350d4fde5777cf0e8256857442 as dakara_check

FROM ghcr.io/odrling/chimera:x86_64@sha256:41b1908517a31791616a7a5d6b89351d227331816ad5bbd2755ace75f40d9cee AS builder 

RUN apk add chimera-repo-contrib

ARG ARCH
ARG GOARCH
ARG SYSROOT

COPY --from=dakara_check /usr/local $SYSROOT/usr/

COPY . /karaberus

RUN cd /karaberus && ARCH=$ARCH GOARCH=$GOARCH SYSROOT=$SYSROOT ci/build.sh

FROM ghcr.io/odrling/chimera@sha256:26016bae5fc810a109cbadc449be1926594a5819223348dc9a41af20c30dc3c6

COPY --from=builder /karaberus/build/karaberus /
ENV KARABERUS_LISTEN_ADDR=":8888"
EXPOSE 8888
ENTRYPOINT ["/karaberus"]
