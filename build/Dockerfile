FROM alpine:latest

RUN apk update && \
    apk add git && \
    addgroup -g 1000 abc && \
    adduser -u 1000 -G abc -h /home/abc -D abc

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY bld/grafana-exporter-$TARGETOS-$TARGETARCH /app/grafana-exporter
COPY init/grafana-exporter.sh .

USER abc
ENTRYPOINT ["/bin/sh", "/app/grafana-exporter.sh"]
CMD []