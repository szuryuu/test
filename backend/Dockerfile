# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/kasiraiai ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl
ENV TZ=Asia/Jakarta

WORKDIR /app

# Binary
COPY --from=builder /app/bin/kasiraiai .

# Migrations (mounted relative to working dir — app reads "file://migrations")
COPY --from=builder /app/migrations ./migrations

# Non-root user
RUN adduser -D -H -h /app appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=3 \
  CMD curl -sf http://localhost:8080/health || exit 1

ENTRYPOINT ["./kasiraiai"]
