# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-voicevox .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/mcp-voicevox .

# Create temp directory
RUN mkdir -p /tmp/mcp-voicevox && \
    chown -R appuser:appgroup /tmp/mcp-voicevox /app

# Switch to non-root user
USER appuser

# Expose port (for server mode)
EXPOSE 8080

# Default command
CMD ["./mcp-voicevox", "server"]
