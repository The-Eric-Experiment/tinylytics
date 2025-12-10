# syntax=docker/dockerfile:1
# Build stage
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

WORKDIR /build/server

# Copy backend source
COPY ./server .

# Download Go modules and build
RUN go mod download && go build -o tinylytics .

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy from the new path
COPY --from=builder /build/server/tinylytics .
COPY --from=builder /build/server/static ./static
COPY --from=builder /build/server/templates ./templates

EXPOSE 8080

CMD ["./tinylytics"]