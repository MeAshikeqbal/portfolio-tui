FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/portfolio-tui .

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates libcap && \
    addgroup -S portfolio && \
    adduser -S -G portfolio -h /app portfolio && \
    mkdir -p /app/logs /app/.ssh && \
    chown -R portfolio:portfolio /app

COPY --from=builder /out/portfolio-tui /usr/local/bin/portfolio-tui
COPY --chown=portfolio:portfolio config.example.yaml /app/config.example.yaml
COPY --chown=portfolio:portfolio docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN chmod +x /usr/local/bin/docker-entrypoint.sh && \
    setcap 'cap_net_bind_service=+ep' /usr/local/bin/portfolio-tui

USER portfolio

EXPOSE 23234

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["serve"]
