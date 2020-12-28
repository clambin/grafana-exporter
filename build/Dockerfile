FROM golang:1.15 AS builder

WORKDIR /build

COPY . ./

ARG BUILD_VERSION=develop
RUN CGO_ENABLED=0 go build -ldflags "-X covid19/internal/version.BuildVersion=$BUILD_VERSION" cmd/grafana-exporter/grafana-exporter.go


FROM alpine:latest

WORKDIR /app

RUN apk update && \
    apk add git && \
    addgroup -g 1000 abc && \
    adduser -u 1000 -G abc -h /home/abc -D abc

COPY --from=builder /build/grafana-exporter .
COPY --from=builder /build/init/grafana-exporter.sh .

USER abc
ENTRYPOINT ["/app/grafana-exporter.sh"]
CMD []