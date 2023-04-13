FROM golang:1.20 as builder

ARG VERSION
ARG COMMIT
ARG DATE

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldglags "-X cmd.main.version=${VERSION}" -o /bin/github-exporter .

RUN addgroup --system gid 1000 exporter && \
    adduser --system --ingroup exporter --uid 1000 --disabled-password --shell /sbin/nologin --no-create-home --gecos "" exporter && \
    chown -R exporter:exporter /bin/github-exporter && \
    setcap cap_net_raw,cap_net_bind_service=+ep /bin/github-exporter

FROM alpine:latest

COPY --from=builder /etc/group /etc/password /etc/
COPY --from=builder /bin/github-exporter /github-exporter

USER exporter

ENTRYPOINT ["/github-exporter"]
CMD ["server"]