# ==============================================================================
# GitHub Repo Transfer Engine — Multi-Stage Production Dockerfile (Go + Resty)
# ==============================================================================
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically linked binary with stripped debug symbols
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/repo-transfer-server .

# ------------------------------------------------------------------------------
# Production Runtime Stage
# ------------------------------------------------------------------------------
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /bin/repo-transfer-server /app/repo-transfer-server

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=20s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ || exit 1

ENTRYPOINT ["/app/repo-transfer-server"]
