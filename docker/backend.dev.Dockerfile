# Development Dockerfile for Backend
FROM golang:1.24-alpine

# Set environment to bypass SSL verification for Go modules
ENV GOINSECURE=*
ENV GOPRIVATE=none

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Run with go run, using module mode (no vendor)
CMD ["go", "run", "-mod=mod", "main.go"]
