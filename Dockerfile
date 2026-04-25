# Multi-stage build for optimal size
FROM golang:1.21-alpine AS builder

# Install git (needed for go mod)
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o rootme-server cmd/server/main.go

# Final stage - minimal runtime
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S app && \
    adduser -S app -u 1001 -G app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/rootme-server .

# Change ownership
RUN chown app:app rootme-server

# Switch to non-root user
USER app

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Run the server
CMD ["./rootme-server"]