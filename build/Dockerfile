FROM golang:1.15 AS builder

WORKDIR /build

COPY . ./

ARG BUILD_VERSION=develop
RUN CGO_ENABLED=0 go build -ldflags "-X covid19/internal/version.BuildVersion=$BUILD_VERSION" cmd/grafana_exporter/grafana_exporter.go


FROM alpine:latest

WORKDIR /app

RUN RUN apk update && \
    RUN apk add git
RUN addgroup -g 1000 abc && \
    adduser --uid 1000 abc abc


COPY --from=builder /build/grafana_exporter .
COPY --from=builder /build/init/grafana_exporter.sh .

USER abc
ENTRYPOINT ["/app/grafana_exporter.sh"]
CMD []