FROM ghcr.io/odrling/chimera:cross AS builder 

ARG CHOST

COPY . /karaberus
RUN cd /karaberus && CHOST=$CHOST ci/build.sh

FROM ghcr.io/odrling/chimera

COPY --from=builder /image /
ENV KARABERUS_LISTEN_ADDR=":8888"
EXPOSE 8888
ENTRYPOINT ["/karaberus"]
