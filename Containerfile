FROM ghcr.io/japan7/dakara_check:master AS builder

RUN apk add go pkgconf clang llvm lld

COPY . /karaberus

RUN cd /karaberus && \
    export CGO_ENABLED=1 && \
    export CC=clang && \
    go build -buildmode=pie -ldflags '-linkmode=external -s' -o build/ .

FROM alpine

COPY --from=builder /karaberus/build/karaberus /
ENV KARABERUS_LISTEN_ADDR="0.0.0.0:8888"
EXPOSE 8888
ENTRYPOINT ["/karaberus"]
