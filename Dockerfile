# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /opencode ./cmd/opencode

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata git bash

# Copy binary from builder
COPY --from=builder /opencode /usr/local/bin/opencode

# Copy documentation
COPY README.md /app/

# Create non-root user
RUN adduser -D -u 1000 opencode
USER opencode

# Set environment
ENV HOME=/home/opencode
ENV OPENCODE_CONFIG_DIR=/home/opencode/.config/opencode
ENV OPENCODE_DATA_DIR=/home/opencode/.local/share/opencode

# Create directories
RUN mkdir -p ${OPENCODE_CONFIG_DIR} ${OPENCODE_DATA_DIR}

# Set working directory
WORKDIR /workspace

# Entry point
ENTRYPOINT ["opencode"]
CMD ["--help"]
