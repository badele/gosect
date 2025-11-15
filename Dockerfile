# Configuration variables
ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.22

# Build stage
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./

# Compile application
RUN CGO_ENABLED=0 GOOS=linux go build -o gosect .

# Runtime stage
ARG ALPINE_VERSION
FROM alpine:${ALPINE_VERSION}

WORKDIR /work

# Copy binary from builder
COPY --from=builder /app/gosect /usr/local/bin/gosect

# Set entrypoint
ENTRYPOINT ["gosect"]

# Default arguments (can be overridden)
CMD []
