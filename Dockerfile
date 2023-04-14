FROM golang:1.20-alpine as builder

WORKDIR /app

RUN apk add --no-cache git libcap ca-certificates && \
    update-ca-certificates 2>/dev/null || true

COPY github-exporter /bin/github-exporter

RUN addgroup --system --gid 1000 exporter && \
    adduser --system --ingroup exporter --uid 1000 --disabled-password --shell /sbin/nologin --no-create-home --gecos "" exporter && \
    chown -R exporter:exporter /bin/github-exporter && \
    setcap cap_net_raw,cap_net_bind_service=+ep /bin/github-exporter

FROM alpine:latest

COPY --from=builder /etc/group /etc/passwd /etc/
COPY --from=builder /bin/github-exporter /github-exporter

USER exporter

ENTRYPOINT ["/github-exporter"]
CMD ["server"]