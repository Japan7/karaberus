FROM ghcr.io/japan7/dakara_check:master AS builder

RUN apk add go pkgconf clang llvm lld

COPY . /karaberus

RUN cd /karaberus && CGO_ENABLED=1 CC=clang go build -o build/ .

FROM alpine

COPY --from=builder /karaberus/build/karaberus /
ENV KARABERUS_LISTEN_ADDR="0.0.0.0:8888"
EXPOSE 8888
ENTRYPOINT ["/karaberus"]
