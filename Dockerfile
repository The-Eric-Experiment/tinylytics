# syntax=docker/dockerfile:1
FROM golang:1.25.5-alpine

# Install build dependencies for CGO and DuckDB
RUN apk add --no-cache \
    git \
    build-base \
    gcc \
    g++ \
    musl-dev

# Set CGO environment variables for DuckDB
ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /app

# Copy go mod files
COPY ./server/go.mod .
COPY ./server/go.sum .

# Download dependencies
RUN go mod download

# Copy client build
COPY ./client/build ./client

# Copy server source
COPY ./server .

# Build with CGO enabled (required for DuckDB)
RUN go build -o tinylytics .

EXPOSE 8080

CMD [ "./tinylytics" ]