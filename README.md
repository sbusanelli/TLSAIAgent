# TLS Hot Reload AI Agent (Go)

A production-ready TLS certificate hot-reload agent with graceful shutdown, feature flags, and comprehensive pre-commit hooks for code quality.

## Quick Start

### Run
```bash
go mod tidy
go run main.go
```

The server runs on `https://localhost:8443`

### Test
Replace `certs/server.crt` and `certs/server.key` with a new cert.
New TLS connections will immediately use the new cert without restart.

## Development Setup

### Prerequisites
- Go 1.22+
- Python 3.7+ (for pre-commit)
- golangci-lint
- git

### Quick Setup
```bash
# Install development dependencies
make install-hooks

# Run all checks
make check
```

Or manually:
```bash
pip install pre-commit
./setup-pre-commit-hooks.sh
```

### Development Commands
```bash
make help              # Show all commands
make build             # Build binary
make run               # Run the agent
make test              # Run tests
make lint              # Run linter
make fmt               # Format code
make run-hooks-all     # Run pre-commit hooks
```

## Features

### Graceful Shutdown
- Signal-based shutdown (SIGTERM/SIGINT)
- Configurable timeout
- Clean resource cleanup

### Certificate Hot-Reload
- File-based certificate watching
- Immediate TLS certificate updates
- No server restart needed

### Feature Flags
Enable/disable features via configuration:
```bash
export TLS_AGENT_FEATURES_GRACEFUL_SHUTDOWN=true
export TLS_AGENT_FEATURES_LOGGING=false
export FEATURES_CONFIG_PATH=features.yaml
./tls-agent
```

See [FEATURES.md](FEATURES.md) for details.

### Code Quality
Pre-commit hooks ensure:
- ✅ Go code formatting (gofmt, gofumpt)
- ✅ Linting (golangci-lint, revive, go vet)
- ✅ Security scanning (gosec, detect-secrets)
- ✅ Tests pass (go test -race)
- ✅ Compilation succeeds (go build)
- ✅ Dependencies are tidy (go mod tidy)

See [PRE_COMMIT_SETUP.md](PRE_COMMIT_SETUP.md) for setup instructions.

## Documentation

- [FEATURES.md](FEATURES.md) - Feature flags configuration
- [FEATURE_FLAGS_IMPLEMENTATION.md](FEATURE_FLAGS_IMPLEMENTATION.md) - Implementation details
- [PRE_COMMIT_SETUP.md](PRE_COMMIT_SETUP.md) - Pre-commit hooks setup guide
- [PRE_COMMIT_QUICK_REFERENCE.md](PRE_COMMIT_QUICK_REFERENCE.md) - Quick commands reference

## Project Structure
```
.
├── main.go                              # Application entry point
├── go.mod / go.sum                      # Dependencies
├── Makefile                             # Development commands
├── .pre-commit-config.yaml              # Hook configuration
├── .golangci.yaml                       # Linter configuration
├── setup-pre-commit-hooks.sh            # Hook installation script
├── internal/
│   ├── agent/                           # Certificate watcher
│   ├── features/                        # Feature flags
│   └── tlsstore/                        # TLS certificate store
├── certs/                               # TLS certificates
└── .github/workflows/pre-commit.yml     # CI/CD configuration
```

## License

MIT
