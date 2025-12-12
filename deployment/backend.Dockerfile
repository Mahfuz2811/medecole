# ---------- Builder ----------
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS builder

# Install build dependencies required for CGO and database drivers
RUN apk add --no-cache git gcc g++ make musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled (required for MySQL/Redis drivers)
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o main .

# ---------- Runtime ----------
FROM alpine:3.20

# Only install lightweight runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    wget \
    libc6-compat

# Set application timezone
ENV TZ=Asia/Dhaka

# Create non-root user
RUN adduser -D -h /app appuser

WORKDIR /app

# Copy binary compiled in builder
COPY --from=builder --chown=appuser:appuser /build/main /app/main

# Create logs directory with correct permissions
RUN mkdir -p /app/logs && chown -R appuser:appuser /app/logs

# Switch to non-root user
USER appuser

EXPOSE 8080

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

CMD ["/app/main"]
