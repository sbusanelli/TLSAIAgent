# Multi-stage Dockerfile for TLSAIAgent with enterprise-grade security and optimization
# Build Stage
FROM golang:1.22-alpine AS builder

# Set build arguments
ARG VERSION=dev
ARG BUILD_DATE
ARG COMMIT_SHA
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download && \
    go mod verify

# Copy source code
COPY . .

# Build application with security and optimization flags
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS:-linux} \
    GOARCH=${TARGETARCH:-amd64} \
    go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.commitSHA=${COMMIT_SHA}" \
    -a -installsuffix cgo \
    -o tlsai-agent \
    ./...

# Security Scan Stage
FROM builder AS security-scan

# Install security scanning tools
RUN apk add --no-cache \
    curl \
    wget \
    && wget -O /tmp/trivy https://github.com/aquasecurity/trivy/releases/latest/download/trivy_${TARGETOS:-linux}_${TARGETARCH:-amd64} \
    && chmod +x /tmp/trivy \
    && mv /tmp/trivy /usr/local/bin/

# Run security scan on the built binary
RUN /usr/local/bin/trivy fs --scanners vuln,secret --format json /app > /tmp/security-scan.json || true

# Production Stage
FROM scratch AS production

# Import CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Import user and group from builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Create non-root user
USER 65534:65534

# Copy binary from builder
COPY --from=builder /app/tlsai-agent /tlsai-agent

# Copy security scan report
COPY --from=security-scan /tmp/security-scan.json /security-scan.json

# Set entrypoint
ENTRYPOINT ["/tlsai-agent"]

# Default command
CMD ["--help"]

# Labels for enterprise metadata
LABEL org.opencontainers.image.title="TLSAIAgent" \
      org.opencontainers.image.description="Enterprise-grade TLS Certificate Management Agent" \
      org.opencontainers.image.vendor="Sreedhar Busanelli" \
      org.opencontainers.image.version=${VERSION} \
      org.opencontainers.image.created=${BUILD_DATE} \
      org.opencontainers.image.revision=${COMMIT_SHA} \
      org.opencontainers.image.source="https://github.com/sbusanelli/TLSAIAgent" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.documentation="https://github.com/sbusanelli/TLSAIAgent/blob/main/README.md" \
      security.scan="/security-scan.json" \
      maintainer="sbusanelli@example.com"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/tlsai-agent", "--health-check"] || exit 1

# Runtime Stage (for development and debugging)
FROM golang:1.22-alpine AS runtime

# Install runtime dependencies
RUN apk add --no-cache \
    curl \
    ca-certificates \
    tzdata \
    dumb-init \
    && update-ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/tlsai-agent /tlsai-agent

# Copy configuration files
COPY --from=builder /app/features.example.json /etc/tlsai-agent/features.json
COPY --from=builder /app/features.example.yaml /etc/tlsai-agent/features.yaml

# Create non-root user
RUN addgroup -g 1001 -S tlsai && \
    adduser -u 1001 -S tlsai -G tlsai

# Change ownership
RUN chown -R tlsai:tlsai /app /etc/tlsai-agent

# Switch to non-root user
USER tlsai

# Expose ports
EXPOSE 8080 8443

# Set entrypoint with dumb-init for proper signal handling
ENTRYPOINT ["dumb-init", "--"]

# Default command
CMD ["/tlsai-agent", "--config", "/etc/tlsai-agent/features.yaml"]

# Labels for runtime stage
LABEL org.opencontainers.image.title="TLSAIAgent Runtime" \
      org.opencontainers.image.description="TLSAIAgent runtime environment with debugging support" \
      org.opencontainers.image.vendor="Sreedhar Busanelli" \
      org.opencontainers.image.version=${VERSION} \
      org.opencontainers.image.created=${BUILD_DATE} \
      org.opencontainers.image.revision=${COMMIT_SHA}

# Health check for runtime
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Development Stage
FROM runtime AS development

# Install development tools
RUN apk add --no-cache \
    git \
    vim \
    htop \
    strace \
    lsof \
    tcpdump

# Copy source code for development
COPY . .

# Set development environment variables
ENV GO_ENV=development
ENV LOG_LEVEL=debug
ENV HOT_RELOAD=true

# Default development command
CMD ["go", "run", "main.go", "--config", "/etc/tlsai-agent/features.yaml", "--log-level", "debug"]

# Labels for development stage
LABEL org.opencontainers.image.title="TLSAIAgent Development" \
      org.opencontainers.image.description="TLSAIAgent development environment with hot reload" \
      org.opencontainers.image.vendor="Sreedhar Busanelli" \
      org.opencontainers.image.version=${VERSION} \
      org.opencontainers.image.created=${BUILD_DATE} \
      org.opencontainers.image.revision=${COMMIT_SHA}
