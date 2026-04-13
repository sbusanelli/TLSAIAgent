# 🔐 TLSAIAgent

Production-ready TLS certificate hot-reload agent with graceful shutdown and feature flags. Go-based service for automatic TLS certificate rotation with zero-downtime updates.

## 🚀 Features

### Core Functionality
- **🔄 Automatic Certificate Hot-Reload**: Monitors and reloads TLS certificates without service interruption
- **⚡ Zero-Downtime Updates**: Seamless certificate rotation with graceful connection handling
- **🛡️ Secure TLS Configuration**: Enforces TLS 1.2+ with modern security standards
- **📋 Feature Flags**: Granular control over agent functionality via configuration

### Advanced Features
- **🎯 Graceful Shutdown**: Clean service termination with configurable timeouts
- **📊 Comprehensive Logging**: Detailed operational logs with configurable verbosity
- **⚙️ Flexible Configuration**: Support for YAML, JSON, and environment variable configuration
- **🔍 Certificate Monitoring**: Real-time file system monitoring for certificate changes

## 📋 Prerequisites

- **Go 1.22+**: Required for building and running the agent
- **TLS Certificates**: Valid certificate and key files in `certs/` directory
- **Linux/macOS/Windows**: Cross-platform support

## 🛠️ Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/sbusanelli/TLSAIAgent.git
cd TLSAIAgent

# Build the application
go build -o tlsai-agent main.go

# Or use make
make build
```

### Basic Usage

```bash
# Start with default configuration
./tlsai-agent

# Start with custom feature configuration
FEATURES_CONFIG_PATH=config/features.yaml ./tlsai-agent

# Start with environment variables
CERTIFICATE_WATCHER=true GRACEFUL_SHUTDOWN=true LOGGING=true ./tlsai-agent
```

### Certificate Setup

```bash
# Create certs directory
mkdir -p certs

# Place your certificates
cp your-server.crt certs/server.crt
cp your-server.key certs/server.key

# Or generate self-signed certificates for testing
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -subj "/CN=localhost"
```

## ⚙️ Configuration

### Feature Flags

Configure agent behavior using feature flags:

```yaml
# config/features.yaml
certificate_watcher: true    # Enable certificate hot-reload
graceful_shutdown: true     # Enable graceful shutdown
logging: true              # Enable detailed logging
shutdown_timeout: 30       # Shutdown timeout in seconds
agent_shutdown_timeout: 10  # Agent shutdown timeout in seconds
```

### Environment Variables

```bash
# Override configuration with environment variables
export CERTIFICATE_WATCHER=true
export GRACEFUL_SHUTDOWN=true
export LOGGING=true
export SHUTDOWN_TIMEOUT=30
export AGENT_SHUTDOWN_TIMEOUT=10
export FEATURES_CONFIG_PATH=config/features.yaml
```

### JSON Configuration

```json
{
  "certificate_watcher": true,
  "graceful_shutdown": true,
  "logging": true,
  "shutdown_timeout": 30,
  "agent_shutdown_timeout": 10
}
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    TLSAIAgent Architecture                    │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   HTTP      │  │   TLS       │  │ Certificate │         │
│  │   Server    │  │   Config    │  │   Store     │         │
│  │  (8443)     │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│         │               │               │                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Feature   │  │  File Watch │  │   Graceful  │         │
│  │   Flags     │  │   Agent     │  │  Shutdown   │         │
│  │  Manager    │  │             │  │  Handler    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────────┘
```

## 📊 Usage Examples

### Development Mode

```bash
# Enable all features for development
export CERTIFICATE_WATCHER=true
export GRACEFUL_SHUTDOWN=true
export LOGGING=true

./tlsai-agent
```

### Production Mode

```bash
# Production configuration
./tlsai-agent \
  --features-config=config/production.yaml \
  --cert-dir=/etc/ssl/certs
```

### Docker Deployment

```bash
# Build Docker image
docker build -t tlsai-agent .

# Run with Docker
docker run -d \
  --name tlsai-agent \
  -p 8443:8443 \
  -v $(pwd)/certs:/app/certs \
  -v $(pwd)/config:/app/config \
  tlsai-agent
```

## 🔧 Development

### Project Structure

```
.
├── main.go                              # Application entry point
├── go.mod / go.sum                      # Dependencies
├── Makefile                             # Development commands
├── .pre-commit-config.yaml              # Pre-commit hooks
├── .golangci.yaml                       # Linter configuration
├── internal/
│   ├── agent/                           # Certificate watcher
│   ├── features/                        # Feature flags
│   └── tlsstore/                        # TLS certificate store
├── certs/                               # TLS certificates
├── config/                              # Configuration files
└── .github/workflows/                   # CI/CD pipelines
```

### Development Commands

```bash
# Install dependencies
go mod tidy

# Run tests
go test -v ./...

# Run with race detection
go test -race -v ./...

# Format code
go fmt ./...

# Run linter
golangci-lint run

# Build for development
go build -o tlsai-agent main.go

# Build for production
go build -ldflags="-s -w" -o tlsai-agent main.go
```

### Pre-commit Hooks

The project uses pre-commit hooks to ensure code quality:

```bash
# Install pre-commit hooks
make setup-hooks

# Or manually
pre-commit install

# Run hooks manually
pre-commit run --all-files
```

## 🐳 Docker Support

### Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o tlsai-agent main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/tlsai-agent .
COPY --from=builder /app/certs ./certs
EXPOSE 8443
CMD ["./tlsai-agent"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  tlsai-agent:
    build: .
    ports:
      - "8443:8443"
    volumes:
      - ./certs:/app/certs
      - ./config:/app/config
    environment:
      - CERTIFICATE_WATCHER=true
      - GRACEFUL_SHUTDOWN=true
      - LOGGING=true
    restart: unless-stopped
```

## 🔒 Security Features

- **TLS 1.2+ Enforcement**: Minimum TLS version 1.2 for secure connections
- **Certificate Validation**: Automatic certificate validation before loading
- **Secure File Watching**: Safe file system monitoring with proper permissions
- **Graceful Shutdown**: Secure termination without data loss

## 📝 Logging

The agent provides comprehensive logging:

```bash
# Enable verbose logging
export LOGGING=true

# Application logs
./tlsai-agent

# Log output example
🎨 TLS Agent server running on https://localhost:8443
   Press Ctrl+C to gracefully shutdown

Certificate watcher agent started
Server shutdown complete
Agent stopped gracefully
TLS Agent shutdown complete
```

## 🧪 Testing

```bash
# Run all tests
go test -v ./...

# Run specific test
go test -v ./internal/agent

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 📚 Documentation

- [FEATURES.md](FEATURES.md) - Feature flags configuration
- [FEATURE_FLAGS_IMPLEMENTATION.md](FEATURE_FLAGS_IMPLEMENTATION.md) - Implementation details
- [PRE_COMMIT_SETUP.md](PRE_COMMIT_SETUP.md) - Pre-commit hooks setup
- [docs/deployment-guide.md](docs/deployment-guide.md) - Deployment guide
- [docs/ENVIRONMENT_PROMOTION.md](docs/ENVIRONMENT_PROMOTION.md) - Environment promotion

## 🔄 CI/CD

The project includes comprehensive GitHub Actions workflows:

- **dependency-update.yml** - Automated dependency management
- **security-scan.yml** - Security vulnerability scanning
- **deploy.yml** - Deployment pipeline
- **pre-commit.yml** - Code quality checks

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/sbusanelli/TLSAIAgent/issues)
- **Documentation**: [GitHub Wiki](https://github.com/sbusanelli/TLSAIAgent/wiki)
- **Email**: sbusanelli@example.com

---

**🔐 TLSAIAgent** - Secure, reliable, and production-ready certificate management.
