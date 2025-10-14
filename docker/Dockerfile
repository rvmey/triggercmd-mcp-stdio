# Multi-stage build for TRIGGERcmd MCP Server
FROM golang:1.23-alpine AS builder

# Install git (needed for Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the binary for Linux AMD64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o triggercmd-mcp .

# Final stage - minimal runtime image
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 -S triggercmd && \
    adduser -u 1000 -S triggercmd -G triggercmd

# Create directory for token file
RUN mkdir -p /home/triggercmd/.TRIGGERcmdData && \
    chown -R triggercmd:triggercmd /home/triggercmd

# Set working directory
WORKDIR /home/triggercmd

# Copy binary from builder stage
COPY --from=builder /app/triggercmd-mcp /usr/local/bin/triggercmd-mcp

# Make binary executable
RUN chmod +x /usr/local/bin/triggercmd-mcp

# Switch to non-root user
USER triggercmd

# Expose no ports (stdio only)

# Set environment variables
ENV HOME=/home/triggercmd

# Health check (optional - checks if binary exists and is executable)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD test -x /usr/local/bin/triggercmd-mcp || exit 1

# Run the MCP server
ENTRYPOINT ["/usr/local/bin/triggercmd-mcp"]