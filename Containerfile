FROM ghcr.io/odrling/chimera:cross AS builder

ARG CHOST

COPY . /karaberus
RUN cd /karaberus && CHOST=$CHOST ci/build.sh

FROM ghcr.io/odrling/chimera

COPY --from=builder /image /
ENV KARABERUS_UI_DIST_DIR="/usr/local/share/karaberus/ui_dist"
EXPOSE 8888
ENTRYPOINT ["karaberus"]
