# syntax=docker/dockerfile:1

# Build stage - use Debian-based Go image for glibc compatibility with DuckDB
FROM golang:1.25.5-bookworm AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Enable CGO for DuckDB (REQUIRED)
ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /app

# Copy go dependencies
COPY ./server/go.mod .
COPY ./server/go.sum .

# Download Go modules
RUN go mod download

# Copy frontend build
COPY ./client/build ./client

# Copy backend source
COPY ./server .

# Build the application
RUN go build -o tinylytics .

# Runtime stage - minimal Debian image (glibc compatible)
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary and client from builder
COPY --from=builder /app/tinylytics .
COPY --from=builder /app/client ./client

EXPOSE 8080

CMD ["./tinylytics"]