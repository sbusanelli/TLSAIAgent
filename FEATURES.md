# Feature Flags Configuration Guide

## Overview

The TLS Agent now supports flexible feature configuration through feature flags. Users can enable/disable features and customize timeouts without modifying code. Configuration can be loaded from YAML files, JSON files, or environment variables.

## Quick Start

### Default Configuration

By default, the TLS Agent runs with these settings:

```yaml
graceful_shutdown: true           # ✅ Enabled
certificate_watcher: true         # ✅ Enabled
periodic_cert_check: true        # ✅ Enabled
debounce_file_changes: true       # ✅ Enabled
logging: true                     # ✅ Enabled
metrics_collection: false         # ❌ Disabled (future feature)
health_check: false              # ❌ Disabled (future feature)
shutdown_timeout: 10              # 10 seconds
agent_shutdown_timeout: 5         # 5 seconds
cert_watch_interval: 30           # 30 seconds
debounce_interval: 2000           # 2 seconds (2000ms)
cert_expiry_warning: 7            # 7 days
```

### Configuration Methods

#### 1. YAML Configuration File

Create a `features.yaml` file:

```yaml
graceful_shutdown: true
certificate_watcher: true
periodic_cert_check: true
debounce_file_changes: true
logging: true
metrics_collection: false
health_check: false
shutdown_timeout: 10
agent_shutdown_timeout: 5
cert_watch_interval: 30
debounce_interval: 2000
cert_expiry_warning: 7
```

Load it:

```bash
export FEATURES_CONFIG_PATH=/path/to/features.yaml
./tls-agent
```

#### 2. JSON Configuration File

Create a `features.json` file:

```json
{
  "graceful_shutdown": true,
  "certificate_watcher": true,
  "periodic_cert_check": true,
  "debounce_file_changes": true,
  "logging": true,
  "metrics_collection": false,
  "health_check": false,
  "shutdown_timeout": 10,
  "agent_shutdown_timeout": 5,
  "cert_watch_interval": 30,
  "debounce_interval": 2000,
  "cert_expiry_warning": 7
}
```

Load it:

```bash
export FEATURES_CONFIG_PATH=/path/to/features.json
./tls-agent
```

#### 3. Environment Variables

Override specific features using environment variables:

```bash
# Boolean features
export TLS_AGENT_FEATURES_GRACEFUL_SHUTDOWN=true
export TLS_AGENT_FEATURES_CERTIFICATE_WATCHER=false
export TLS_AGENT_FEATURES_LOGGING=true

# Integer configurations
export TLS_AGENT_FEATURES_SHUTDOWN_TIMEOUT=15
export TLS_AGENT_FEATURES_AGENT_SHUTDOWN_TIMEOUT=8
export TLS_AGENT_FEATURES_CERT_WATCH_INTERVAL=60

./tls-agent
```

## Features Explained

### Boolean Features

#### `graceful_shutdown` (default: `true`)

When enabled, the server gracefully shuts down when receiving SIGTERM or SIGINT signals. It:
- Stops accepting new connections
- Waits for existing connections to complete
- Gracefully stops the certificate watcher agent
- Uses configured timeout to force shutdown if needed

**When to disable:** In resource-constrained environments or for development where immediate shutdown is acceptable.

#### `certificate_watcher` (default: `true`)

When enabled, starts the certificate file watcher agent that:
- Monitors certificate and key file changes
- Automatically reloads certificates when they change
- Enables hot-reload of TLS certificates without restarting the server
- Uses periodic checks as fallback

**When to disable:** If certificates are manually rotated and server restarts are acceptable.

#### `periodic_cert_check` (default: `true`)

When enabled, performs periodic certificate checks:
- Detects certificate expiry dates
- Warns before certificates expire (based on `cert_expiry_warning` setting)
- Provides redundancy for file-based certificate watching

**When to disable:** If using external monitoring systems for certificate expiry alerts.

#### `debounce_file_changes` (default: `true`)

When enabled, debounces rapid certificate file changes:
- Prevents reload storms if files are updated multiple times rapidly
- Uses `debounce_interval` to batch multiple changes into one reload
- Protects server stability during bulk certificate updates

**When to disable:** For immediate reload on any file change (not recommended).

#### `logging` (default: `true`)

When enabled, outputs detailed logs including:
- Feature configuration at startup
- Certificate reload events
- Shutdown process details
- Performance metrics

**When to disable:** In production for cleaner output or when using centralized logging.

#### `metrics_collection` (default: `false`)

When enabled, collects performance metrics (future feature):
- Currently a placeholder for future enhancement
- Will track response times, certificate reloads, etc.

#### `health_check` (default: `false`)

When enabled, provides a health check endpoint (future feature):
- Currently a placeholder for future enhancement
- Will provide `/health` endpoint for Kubernetes/load balancer checks

## Integer Configurations

### `shutdown_timeout` (default: `10` seconds)

Maximum time to wait for HTTP server graceful shutdown. If exceeded, the server forcibly closes connections.

**Examples:**
- `5` - Quick shutdown (may interrupt connections)
- `10` - Balanced (recommended)
- `30` - Long shutdown (waits for slow clients)

### `agent_shutdown_timeout` (default: `5` seconds)

Maximum time to wait for the certificate watcher agent to stop cleanly.

**Examples:**
- `2` - Quick cleanup
- `5` - Balanced (recommended)
- `10` - Extended cleanup

### `cert_watch_interval` (default: `30` seconds)

Interval between periodic certificate checks (in seconds).

**Examples:**
- `10` - Frequent checks (higher CPU usage)
- `30` - Balanced (recommended)
- `300` - Infrequent checks (lower CPU usage)

### `debounce_interval` (default: `2000` milliseconds)

Time to wait after a file change before reloading certificate (in milliseconds).

**Examples:**
- `500` - Quick reload (2 file changes batched)
- `2000` - Balanced (recommended)
- `5000` - Delayed reload (many changes batched)

### `cert_expiry_warning` (default: `7` days)

Days before certificate expiry to trigger warning logs.

**Examples:**
- `1` - Only warn when certificate expires tomorrow
- `7` - Standard warning period (recommended)
- `30` - Early warning for planning

## Usage Examples

### Example 1: Production Setup (Minimal Overhead)

```yaml
# features.yaml
graceful_shutdown: true
certificate_watcher: true
periodic_cert_check: false        # External monitoring
debounce_file_changes: true
logging: false                    # Use centralized logging
metrics_collection: false
health_check: false
shutdown_timeout: 10
agent_shutdown_timeout: 5
cert_watch_interval: 60           # Less frequent checks
debounce_interval: 2000
cert_expiry_warning: 14           # Earlier warning
```

### Example 2: Development Setup (Maximum Features)

```yaml
# features.yaml
graceful_shutdown: true
certificate_watcher: true
periodic_cert_check: true
debounce_file_changes: true
logging: true                     # Detailed logs
metrics_collection: true          # For debugging
health_check: true
shutdown_timeout: 5               # Quick shutdown
agent_shutdown_timeout: 2
cert_watch_interval: 10           # Frequent updates
debounce_interval: 500            # Minimal debounce
cert_expiry_warning: 1
```

### Example 3: Minimal Setup (Emergency/Fallback)

```bash
export TLS_AGENT_FEATURES_GRACEFUL_SHUTDOWN=false
export TLS_AGENT_FEATURES_CERTIFICATE_WATCHER=false
export TLS_AGENT_FEATURES_LOGGING=false
./tls-agent
```

### Example 4: Docker Environment

```dockerfile
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o tls-agent

FROM alpine:latest
COPY --from=builder /app/tls-agent /usr/local/bin/
COPY --from=builder /app/features.yaml /etc/tls-agent/
COPY certs/ /etc/tls-agent/certs/

ENV FEATURES_CONFIG_PATH=/etc/tls-agent/features.yaml
ENV TLS_AGENT_FEATURES_LOGGING=true
EXPOSE 8443
CMD ["tls-agent"]
```

### Example 5: Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tls-agent
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: tls-agent
        image: myregistry/tls-agent:latest
        ports:
        - containerPort: 8443
        env:
        - name: FEATURES_CONFIG_PATH
          value: /etc/config/features.yaml
        - name: TLS_AGENT_FEATURES_LOGGING
          value: "true"
        - name: TLS_AGENT_FEATURES_GRACEFUL_SHUTDOWN
          value: "true"
        volumeMounts:
        - name: config
          mountPath: /etc/config
        - name: certs
          mountPath: /etc/certs
      volumes:
      - name: config
        configMap:
          name: tls-agent-config
      - name: certs
        secret:
          secretName: tls-agent-certs
```

## Configuration Priority

When multiple configuration methods are used, they are applied in this order:

1. **Default configuration** - Built-in defaults
2. **Config file** (`features.yaml` or `features.json`) - If `FEATURES_CONFIG_PATH` is set
3. **Environment variables** - Takes highest priority, overrides all others

**Example:** If you set `FEATURES_CONFIG_PATH=features.yaml` and also set `TLS_AGENT_FEATURES_LOGGING=false`, the logging feature will be disabled regardless of the YAML file.

## Troubleshooting

### Features not loading from config file

```bash
# Check if file exists
ls -la /path/to/features.yaml

# Ensure environment variable is set correctly
echo $FEATURES_CONFIG_PATH

# Check logs for errors
./tls-agent 2>&1 | grep -i "feature\|config"
```

### Certificate reloads not happening

```bash
# Enable logging to see reload events
export TLS_AGENT_FEATURES_LOGGING=true

# Check certificate watcher is enabled
export TLS_AGENT_FEATURES_CERTIFICATE_WATCHER=true

# Watch the logs
./tls-agent
```

### Graceful shutdown taking too long

```bash
# Reduce timeouts
export TLS_AGENT_FEATURES_SHUTDOWN_TIMEOUT=5
export TLS_AGENT_FEATURES_AGENT_SHUTDOWN_TIMEOUT=2

./tls-agent
```

## Migration from Hardcoded Configuration

If you previously relied on hardcoded values in the code:

1. **Identify current values** in `main.go`
2. **Create features.yaml** with matching values
3. **Set FEATURES_CONFIG_PATH** to point to the file
4. **Test thoroughly** in staging environment
5. **Deploy** to production

Example migration for shutdown timeout:

```go
// Before (hardcoded)
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

// After (feature-configured)
ctx, cancel := context.WithTimeout(context.Background(), 
  time.Duration(featureConfig.ShutdownTimeout)*time.Second)
```

## Example Files

See the following example files in the project root:

- `features.example.yaml` - YAML configuration example
- `features.example.json` - JSON configuration example

Copy and customize one of these for your deployment!
