FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ENV VERSION=$VERSION

WORKDIR /app/
ADD . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-X main.version=$VERSION" \
    -o grafana-exporter \
    cmd/grafana-exporter/grafana-exporter.go


FROM alpine:latest

RUN apk update && \
    apk add git && \
    addgroup -g 1000 abc && \
    adduser -u 1000 -G abc -h /home/abc -D abc
USER abc

WORKDIR /app
COPY --from=builder /app/grafana-exporter /app/grafana-exporter


ARG TARGETOS
ARG TARGETARCH

ENTRYPOINT ["/app/grafana-exporter"]
CMD []