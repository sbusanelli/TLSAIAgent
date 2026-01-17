# TLS Agent Testing Guide

## 🧪 Overview

This guide covers the comprehensive testing strategy for the TLS Agent, including unit tests, integration tests, benchmarks, and performance testing using modern Go testing practices.

## 🎯 Testing Strategy

### **Testing Pyramid**
```
    Integration Tests (20%)
       ↓
    Unit Tests (70%)
       ↓
    Benchmark Tests (10%)
```

- **Unit Tests**: Fast, isolated tests for individual components
- **Integration Tests**: Tests for component interactions
- **Benchmark Tests**: Performance and memory usage testing

### **Test Coverage Goals**
- **Unit Tests**: 85%+ line coverage
- **Integration Tests**: 80%+ line coverage
- **Benchmarks**: Performance baseline establishment

## 🛠️ Testing Stack

### **Core Testing Tools**
- **Go Testing**: Built-in Go testing framework
- **Testify**: Assertion library for Go
- **Benchmarking**: Go's built-in benchmarking tools
- **Race Detection**: Go's race detector for concurrent code

### **Testing Dependencies**
```go
// No external dependencies required
// Uses Go's built-in testing package
import (
    "testing"
    "time"
    "context"
    "crypto/tls"
)
```

## 📁 Test Structure

```
internal/
├── agent/
│   ├── agent.go
│   └── agent_test.go          # Agent unit tests
├── tlsstore/
│   ├── tlsstore.go
│   └── tlsstore_test.go       # TLS store unit tests
├── features/
│   ├── features.go
│   └── features_test.go       # Feature flag tests
└── ...

integration_test.go         # Integration tests
benchmark_test.go            # Benchmark tests
main_test.go                # Main application tests
```

## 🔬 Unit Testing

### **Agent Tests**
```go
// TestAgentStartStop tests basic agent start and stop functionality
func TestAgentStartStop(t *testing.T) {
    cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
    if err != nil {
        t.Fatalf("Failed to load certificates: %v", err)
    }

    store := tlsstore.New(cert)
    state := agent.NewState(cert)
    agentStopChan := make(chan struct{})
    agentDone := make(chan struct{})

    // Start the agent
    go func() {
        agent.Run(store, state, agentStopChan)
        close(agentDone)
    }()

    // Give agent time to start
    time.Sleep(100 * time.Millisecond)

    // Signal agent to stop
    close(agentStopChan)

    // Wait for agent to stop with timeout
    select {
    case <-agentDone:
        t.Log("Agent stopped successfully")
    case <-time.After(5 * time.Second):
        t.Error("Agent did not stop within timeout")
    }
}
```

### **TLS Store Tests**
```go
// TestLoad tests certificate loading functionality
func TestLoad(t *testing.T) {
    // Test loading valid certificates
    cert, err := Load("certs/server.crt", "certs/server.key")
    if err != nil {
        t.Fatalf("Failed to load certificates: %v", err)
    }

    if cert == nil {
        t.Fatal("Certificate should not be nil")
    }

    if cert.Certificate == nil {
        t.Error("Certificate certificate should not be nil")
    }

    if cert.PrivateKey == nil {
        t.Error("Certificate private key should not be nil")
    }
}
```

### **Feature Flag Tests**
```go
// TestDefaultFeatures verifies the default feature configuration
func TestDefaultFeatures(t *testing.T) {
    features := DefaultFeatures()

    if !features.GracefulShutdown {
        t.Error("GracefulShutdown should be enabled by default")
    }
    if !features.CertificateWatcher {
        t.Error("CertificateWatcher should be enabled by default")
    }
    if !features.Logging {
        t.Error("Logging should be enabled by default")
    }
}
```

## 🔗 Integration Testing

### **Agent-Server Integration**
```go
// TestIntegrationAgentServer tests the complete integration between agent and HTTP server
func TestIntegrationAgentServer(t *testing.T) {
    // Create temporary directory for test certificates
    tempDir := t.TempDir()
    certFile := filepath.Join(tempDir, "server.crt")
    keyFile := filepath.Join(tempDir, "server.key")

    // Create test certificates
    err := createTestCertificates(certFile, keyFile)
    if err != nil {
        t.Fatalf("Failed to create test certificates: %v", err)
    }

    // Load certificates
    cert, err := tlsstore.Load(certFile, keyFile)
    if err != nil {
        t.Fatalf("Failed to load certificates: %v", err)
    }

    store := tlsstore.New(cert)
    state := agent.NewState(cert)
    agentStopChan := make(chan struct{})
    agentDone := make(chan struct{})

    // Start the agent
    go func() {
        agent.Run(store, state, agentStopChan)
        close(agentDone)
    }

    // Give agent time to start
    time.Sleep(100 * time.Millisecond)

    // Create HTTP server
    tlsCfg := &tls.Config{
        GetCertificate: store.GetCertificate,
        MinVersion:     tls.VersionTLS12,
    }

    server := &http.Server{
        Addr:      ":9447",
        TLSConfig: tlsCfg,
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Integration test OK"))
        }),
    }

    // Test HTTP request
    client := &http.Client{
        Timeout: 2 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
        },
    }

    resp, err := client.Get("https://localhost:9447/")
    if err != nil {
        t.Logf("HTTP request failed: %v", err)
    } else {
        resp.Body.Close()
        t.Log("HTTP request succeeded")
    }

    // Cleanup
    // ... shutdown server and agent
}
```

### **Feature Flag Integration**
```go
// TestIntegrationFeatureFlags tests integration with feature flags
func TestIntegrationFeatureFlags(t *testing.T) {
    // Test with feature flags enabled
    features := features.DefaultFeatures()
    features.GracefulShutdown = true
    features.CertificateWatcher = true
    features.Logging = true

    // Create store and state
    store := tlsstore.New(cert)
    state := agent.NewState(cert)

    // Start agent with feature flags
    go func() {
        agent.Run(store, state, agentStopChan)
        close(agentDone)
    }

    // Verify agent is running
    if !state.IsRunning() {
        t.Error("Agent should be running with feature flags enabled")
    }

    // Cleanup
    // ... stop agent
}
```

## ⚡ Benchmark Testing

### **Performance Benchmarks**
```go
// BenchmarkAgentStartup benchmarks agent startup time
func BenchmarkAgentStartup(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        store := tlsstore.New(cert)
        state := agent.NewState(cert)
        agentStopChan := make(chan struct{})
        agentDone := make(chan struct{})

        b.StartTimer()

        // Start the agent
        go func() {
            agent.Run(store, state, agentStopChan)
            close(agentDone)
        }()

        // Give agent time to start
        time.Sleep(10 * time.Millisecond)

        // Signal agent to stop
        close(agentStopChan)

        // Wait for agent to stop
        <-agentDone

        b.StopTimer()
    }
}
```

### **Memory Usage Benchmarks**
```go
// BenchmarkMemoryUsage benchmarks memory usage during operations
func BenchmarkMemoryUsage(b *testing.B) {
    // Measure initial memory usage
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)

    b.ResetTimer()

    // Create multiple stores and states
    stores := make([]*tlsstore.CertificateStore, 100)
    states := make([]*agent.State, 100)

    for i := 0; i < 100; i++ {
        stores[i] = tlsstore.New(cert)
        states[i] = agent.NewState(cert)
    }

    // Perform operations
    for i := 0; i < 1000; i++ {
        for j := 0; j < 100; j++ {
            stores[j].GetCertificate()
            states[j].IncrementCertificateCount()
            states[j].GetCertificateCount()
            states[j].IsRunning()
        }
    }

    b.StopTimer()

    // Measure final memory usage
    runtime.GC()
    runtime.ReadMemStats(&m2)

    memoryUsed := m2.Alloc - m1.Alloc
    b.Logf("Memory used: %d bytes", memoryUsed)
}
```

### **Concurrent Performance**
```go
// BenchmarkConcurrentAccess benchmarks concurrent access to agent state
func BenchmarkConcurrentAccess(b *testing.B) {
    state := agent.NewState(cert)

    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            state.IncrementCertificateCount()
            state.GetCertificateCount()
            state.IsRunning()
            state.GetCertificate()
        }
    })
}
```

## 📊 Test Configuration

### **Makefile Commands**
```makefile
# Test targets
test:
	@echo "🧪 Running tests..."
	@go test -v -race ./...

test-unit:
	@echo "🧪 Running unit tests..."
	@go test -v -race -run "^Test" ./...

test-integration:
	@echo "🔗 Running integration tests..."
	@go test -v -race -run "^TestIntegration" ./...

test-benchmark:
	@echo "⚡ Running benchmark tests..."
	@go test -v -bench=. -benchmem ./...

test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-all: test-unit test-integration test-benchmark test-coverage
```

### **Running Tests**
```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run benchmark tests
make test-benchmark

# Run tests with coverage
make test-coverage

# Run all test suites
make test-all

# Run CI test suite
make test-ci
```

## 📈 Test Coverage

### **Coverage Reports**
- **HTML Report**: `coverage.html`
- **Text Report**: Console output
- **Function Coverage**: Detailed function coverage
- **Line Coverage**: Line-by-line coverage

### **Coverage Thresholds**
- **Statements**: 85%
- **Functions**: 85%
- **Lines**: 85%
- **Branches**: 80%

### **Coverage Commands**
```bash
# Generate coverage report
make test-coverage

# View coverage in browser
open coverage.html

# Check coverage percentage
go tool cover -func=coverage.out
```

## 🔧 Test Utilities

### **Test Helpers**
```go
// Helper function to create test certificates
func createTestCertificates(certFile, keyFile string) error {
    testCert := `-----BEGIN CERTIFICATE-----
MIIDdzCCAn+gAwIBAgI...
-----END CERTIFICATE-----`

    testKey := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQ...
-----END PRIVATE KEY-----`

    err := os.WriteFile(certFile, []byte(testCert), 0644)
    if err != nil {
        return fmt.Errorf("failed to write certificate file: %v", err)
    }

    err = os.WriteFile(keyFile, []byte(testKey), 0644)
    if err != nil {
        return fmt.Errorf("failed to write key file: %v", err)
    }

    return nil
}

// Helper function to create temporary directory
func createTempDir(t *testing.T) string {
    tempDir := t.TempDir()
    return tempDir
}

// Helper function to wait for goroutine
func waitForGoroutine(duration time.Duration) {
    time.Sleep(duration)
}
```

### **Test Assertions**
```go
// Assert certificate is valid
func assertCertificateValid(t *testing.T, cert *tls.Certificate) {
    if cert == nil {
        t.Fatal("Certificate should not be nil")
    }
    if cert.Certificate == nil {
        t.Error("Certificate certificate should not be nil")
    }
    if cert.PrivateKey == nil {
        t.Error("Certificate private key should not be nil")
    }
}

// Assert agent state
func assertAgentState(t *testing.T, state *agent.State, running bool) {
    if state.IsRunning() != running {
        t.Errorf("Agent running state should be %v", running)
    }
}
```

## 🚀 Performance Testing

### **Performance Metrics**
- **Startup Time**: Time to start agent
- **Shutdown Time**: Time to shutdown agent
- **Memory Usage**: Memory consumption during operations
- **Concurrent Performance**: Performance under load
- **TLS Handshake**: TLS certificate retrieval performance

### **Performance Baselines**
```go
// Performance benchmarks
func BenchmarkAgentStartup(b *testing.B) {
    // Baseline: < 10ms startup time
}

func BenchmarkCertificateRetrieval(b *testing.B) {
    // Baseline: < 1ms retrieval time
}

func BenchmarkConcurrentAccess(b *testing.B) {
    // Baseline: Handle 1000 concurrent requests
}
```

### **Memory Testing**
```go
// Memory usage testing
func TestMemoryLeak(t *testing.T) {
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)

    // Perform operations
    // ... test operations

    runtime.GC()
    runtime.ReadMemStats(&m2)

    memoryUsed := m2.Alloc - m1.Alloc
    if memoryUsed > 1024*1024 { // 1MB threshold
        t.Errorf("Memory usage too high: %d bytes", memoryUsed)
    }
}
```

## 📋 Test Checklist

### **Before Writing Tests**
- [ ] Understand the component functionality
- [ ] Identify test scenarios
- [ ] Plan test structure
- [ ] Prepare test data
- [ ] Set up test environment

### **During Test Development**
- [ ] Follow Go testing conventions
- [ ] Use descriptive test names
- [ ] Test both success and failure cases
- [ ] Use table-driven tests for multiple scenarios
- [ ] Clean up resources in teardown

### **After Test Completion**
- [ ] Run tests locally
- [ ] Check coverage reports
- [ ] Verify test reliability
- [ ] Update documentation
- [ ] Review test performance

## 🔍 Test Best Practices

### **Go Testing Best Practices**
1. **Use Table-Driven Tests**: For multiple test cases
2. **Subtests**: Use `t.Run()` for test organization
3. **Helper Functions**: Extract common test logic
4. **Cleanup**: Use `t.Cleanup()` for resource cleanup
5. **Timeouts**: Use `context.WithTimeout` for time-sensitive tests

### **Test Organization**
```go
func TestAgentOperations(t *testing.T) {
    tests := []struct {
        name     string
        setup    func() (*agent.State, *tlsstore.CertificateStore)
        test     func(*testing.T, *agent.State, *tlsstore.CertificateStore)
        cleanup  func(*agent.State, *tlsstore.CertificateStore)
    }{
        {
            name: "start_stop",
            setup: func() (*agent.State, *tlsstore.CertificateStore) {
                // Setup test
                return state, store
            },
            test: func(t *testing.T, state *agent.State, store *tlsstore.CertificateStore) {
                // Test logic
            },
            cleanup: func(state *agent.State, store *tlsstore.CertificateStore) {
                // Cleanup
            },
        },
        // More test cases...
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            state, store := tc.setup()
            defer tc.cleanup(state, store)
            tc.test(t, state, store)
        })
    }
}
```

### **Error Testing**
```go
func TestLoadInvalidFiles(t *testing.T) {
    // Test non-existent files
    _, err := Load("nonexistent.crt", "nonexistent.key")
    if err == nil {
        t.Error("Loading non-existent files should fail")
    }

    // Test empty files
    tempDir := t.TempDir()
    emptyCert := tempDir + "/empty.crt"
    emptyKey := tempDir + "/empty.key"

    // Create empty files
    err = os.WriteFile(emptyCert, []byte{}, 0644)
    if err != nil {
        t.Fatalf("Failed to create empty cert file: %v", err)
    }

    // Try to load empty files
    _, err = Load(emptyCert, emptyKey)
    if err == nil {
        t.Error("Loading empty certificate files should fail")
    }
}
```

## 📚 Additional Resources

### **Go Testing Documentation**
- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Benchmarking](https://pkg.go.dev/testing#hdr-Benchmarks)
- [Go Race Detector](https://pkg.go/cmd/go/internal/race)

### **Testing Tools**
- [Testify](https://github.com/stretchr/testify)
- [GoMock](https://github.com/golang/mock)
- [Ginkgo](https://onsi.github.io/ginkgo/)
- [Testify](https://github.com/stretchr/testify)

### **Performance Testing**
- [Go Benchmarking](https://golang.org/pkg/testing/)
- [pprof](https://pkg.go/cmd/pprof/)
- [Go Trace](https://golang.org/cmd/trace)

---

**This comprehensive testing strategy ensures the TLS Agent maintains high code quality and reliability!** 🧪
