FROM golang:1.25.6-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/quick-admin ./cmd/api

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata wget \
    && addgroup -S quickadmin \
    && adduser -S -G quickadmin quickadmin

WORKDIR /app
COPY --from=builder /out/quick-admin /app/quick-admin
COPY cmd/api/conf.prod.yaml /app/conf.prod.yaml
COPY cmd/api/zaplogger.prod.yaml /app/zaplogger.prod.yaml

ENV TZ=Asia/Shanghai \
    QUICK_ADMIN_APP_ENV=prod

USER quickadmin
EXPOSE 9009

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://127.0.0.1:9009/health/live || exit 1

ENTRYPOINT ["/app/quick-admin"]
