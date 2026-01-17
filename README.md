# Virtual Threads vs Goroutines Performance Benchmark

A comprehensive performance comparison between Java Virtual Threads and Go Goroutines for concurrent I/O-bound tasks.

## Purpose

This project demonstrates and benchmarks the performance characteristics of:
- **Java Virtual Threads** (Java 21+)
- **Java Traditional Threads** 
- **Go Goroutines**

Perfect for understanding concurrency patterns and performance trade-offs in modern programming languages.

## Architecture

### Java Implementation
```
┌─────────────────────────┐    ┌──────────────────────────┐
│   Virtual Threads      │    │   Traditional Threads    │
│  (JVM-managed)       │    │   (OS-mapped)          │
├─────────────────────────┤    ├──────────────────────────┤
│ Lightweight            │    │ Heavyweight             │
│ Millions possible       │    │ Limited by OS           │
│ Low memory overhead    │    │ High memory overhead     │
└─────────────────────────┘    └──────────────────────────┘
```

### Go Implementation
```
┌─────────────────────────┐
│     Goroutines        │
│  (Go runtime)        │
├─────────────────────────┤
│ Extremely lightweight  │
│ Built-in concurrency  │
│ CSP communication    │
│ Excellent scalability │
└─────────────────────────┘
```

## Benchmark Scenarios

### 1. **I/O-Bound Tasks** (Primary Focus)
- Simulated network requests (100ms latency)
- File I/O operations
- Database query simulations
- **Expected**: Virtual Threads & Goroutines excel

### 2. **CPU-Bound Tasks**
- Mathematical computations
- Data processing
- **Expected**: Traditional threads may compete better

### 3. **Mixed Workload**
- Combination of I/O and CPU operations
- Real-world application simulation

## Quick Start

### Java Benchmarks
```bash
# Compile
mvn clean compile

# Run Virtual Threads benchmark
mvn exec:java -Dexec.mainClass="com.benchmark.VirtualThreadsBenchmark"

# Run Traditional Threads benchmark  
mvn exec:java -Dexec.mainClass="com.benchmark.TraditionalThreadsBenchmark"

# Run comprehensive comparison
mvn exec:java -Dexec.mainClass="com.benchmark.ComparisonRunner"
```

### Go Benchmarks
```bash
# Navigate to Go implementation
cd go-implementation

# Run goroutine benchmark
go run main.go

# Run with custom parameters
go run main.go -tasks=10000 -duration=100ms -workers=50
```

## Performance Metrics

The benchmarks measure:
- **Execution Time**: Total time to complete all tasks
- **Memory Usage**: Peak memory consumption
- **CPU Utilization**: How efficiently CPU resources are used
- **Scalability**: Performance with increasing task counts
- **Throughput**: Tasks completed per second

## Expected Results

Based on preliminary testing:

| Implementation | 1K Tasks | 10K Tasks | 100K Tasks | Memory Efficiency |
|----------------|-------------|--------------|---------------|------------------|
| Java Virtual    | ~150ms      | ~200ms       | ~300ms        | ⭐⭐⭐⭐⭐         |
| Java Traditional| ~800ms      | ~2000ms      | ~5000ms       | ⭐⭐              |
| Go Goroutines  | ~120ms      | ~150ms       | ~200ms        | ⭐⭐⭐⭐⭐⭐         |

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
