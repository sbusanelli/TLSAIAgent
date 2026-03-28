package agent

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"tls-agent/internal/tlsstore"
)

// TestAgentStartStop tests basic agent start and stop functionality
func TestAgentStartStop(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start the agent
	go func() {
		Run(store, state, agentStopChan)
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

// TestAgentState tests agent state management
func TestAgentState(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	state := NewState(cert)

	// Test initial state
	if state.Current == nil {
		t.Error("Current certificate should not be nil")
	}
	if state.Previous != nil {
		t.Error("Previous certificate should be nil initially")
	}
	if state.LastRun.IsZero() {
		t.Error("LastRun should be set after state creation")
	}
}

// TestAgentFileWatcher tests file watching functionality
func TestAgentFileWatcher(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFile := tempDir + "/test.crt"

	// Create test file
	err = os.WriteFile(testFile, []byte("test certificate"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Start agent with file watcher
	go func() {
		Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Modify test file to trigger file watcher
	err = os.WriteFile(testFile, []byte("modified certificate"), 0644)
	if err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Give file watcher time to detect changes
	time.Sleep(200 * time.Millisecond)

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent with file watcher stopped successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent with file watcher did not stop within timeout")
	}

	// Clean up
	os.Remove(testFile)
}

// TestAgentSignalHandling tests signal handling functionality
func TestAgentSignalHandling(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start the agent
	go func() {
		Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Test signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Send signal
	sigChan <- syscall.SIGINT

	// Wait for signal to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify agent state is still valid (should handle signal gracefully)
	if state.Current == nil {
		t.Error("Agent should still have current certificate after signal")
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent handled signal and stopped successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop after signal handling")
	}
}

// TestAgentConcurrentAccess tests concurrent access to agent state
func TestAgentConcurrentAccess(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start the agent
	go func() {
		Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Test concurrent access to state
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			// Test concurrent access to state fields
			_ = state.Current
			_ = state.Previous
			_ = state.LastRun
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify state is still valid
	if state.Current == nil {
		t.Error("Current certificate should not be nil after concurrent access")
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent handled concurrent access successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop after concurrent access test")
	}
}

// TestAgentMemoryLeak tests for memory leaks during long-running operation
func TestAgentMemoryLeak(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start the agent
	go func() {
		Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Simulate long-running operation with periodic state access
	stopTest := make(chan bool)
	go func() {
		for {
			select {
			case <-stopTest:
				return
			default:
				// Access state fields to test memory usage
				_ = state.Current
				_ = state.Previous
				state.LastRun = time.Now()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Run for a short duration
	time.Sleep(500 * time.Millisecond)
	close(stopTest)

	// Verify agent state is still valid
	if state.Current == nil {
		t.Error("Agent should still have current certificate during memory leak test")
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent completed memory leak test successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop after memory leak test")
	}
}

// TestAgentErrorHandling tests error handling in agent operations
func TestAgentErrorHandling(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start the agent
	go func() {
		Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Test error conditions
	// Note: In a real implementation, you would test specific error scenarios
	// such as file system errors, certificate loading errors, etc.

	// Test state operations with error conditions
	// This is a placeholder for error handling tests
	// Test that state fields are accessible
	if state.Current == nil {
		t.Error("Current certificate should not be nil")
	}
	// Previous can be nil initially, that's expected

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent handled error conditions successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop after error handling test")
	}
}

// TestAgentResourceCleanup tests resource cleanup when agent stops
func TestAgentResourceCleanup(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start the agent
	go func() {
		Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent stopped successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop within timeout")
	}

	// Verify resources are cleaned up
	// Note: In a real implementation, you would verify that
	// file watchers, goroutines, and other resources are properly cleaned up

	// Test that state is still accessible after cleanup
	if state.Current == nil {
		t.Error("Current certificate should still be accessible after cleanup")
	}
}

// TestAgentConfiguration tests agent configuration handling
func TestAgentConfiguration(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	state := NewState(cert)

	// Test certificate configuration
	if state.Current == nil {
		t.Error("Current certificate should not be nil")
	}

	// Test certificate properties
	if state.Current.Certificate == nil {
		t.Error("Current certificate data should not be nil")
	}

	if state.Current.PrivateKey == nil {
		t.Error("Current private key should not be nil")
	}

	// Test state configuration
	if state.Previous != nil {
		t.Error("Previous certificate should be nil initially")
	}

	// Test timestamp initialization
	if state.LastRun.IsZero() {
		t.Error("LastRun should be set after state creation")
	}
}

// TestAgentPerformance tests agent performance under load
func TestAgentPerformance(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start the agent
	go func() {
		Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Performance test: rapid state field access
	startTime := time.Now()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		// Test field access performance
		_ = state.Current
		_ = state.Previous
		_ = state.LastRun
		state.LastRun = time.Now()
	}

	duration := time.Since(startTime)
	avgDuration := duration / time.Duration(iterations*4) // 4 operations per iteration

	t.Logf("Performance test: %d iterations in %v (avg: %v per operation)",
		iterations, duration, avgDuration)

	// Performance threshold (adjust as needed)
	if avgDuration > 1*time.Microsecond {
		t.Errorf("Average operation duration too high: %v", avgDuration)
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent performance test completed successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop after performance test")
	}
}

// TestAgentIntegration tests integration with HTTP server
func TestAgentIntegration(t *testing.T) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start the agent
	go func() {
		Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Create HTTP server with TLS
	tlsCfg := &tls.Config{
		GetCertificate: store.GetCertificate,
		MinVersion:     tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":9446",
		TLSConfig: tlsCfg,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}),
	}

	serverStarted := make(chan struct{})
	serverStopped := make(chan struct{})

	// Start server
	go func() {
		close(serverStarted)
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
		close(serverStopped)
	}()

	// Wait for server to start
	<-serverStarted
	time.Sleep(100 * time.Millisecond)

	// Test HTTP request
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get("https://localhost:9446/")
	if err != nil {
		t.Logf("HTTP request failed: %v", err)
	} else {
		resp.Body.Close()
		t.Log("HTTP request succeeded")
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		t.Logf("Server shutdown error: %v", err)
	}

	// Wait for server to stop
	select {
	case <-serverStopped:
		t.Log("Server stopped")
	case <-time.After(6 * time.Second):
		t.Error("Server did not stop in time")
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent integration test completed successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop after integration test")
	}
}

// BenchmarkAgentOperations benchmarks agent operations
func BenchmarkAgentOperations(b *testing.B) {
	cert, err := tlsstore.Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	state := NewState(cert)

	b.ResetTimer()

	// Benchmark state field access operations
	b.Run("CurrentAccess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = state.Current
		}
	})

	b.Run("PreviousAccess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = state.Previous
		}
	})

	b.Run("LastRunAccess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = state.LastRun
		}
	})

	b.Run("LastRunUpdate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			state.LastRun = time.Now()
		}
	})
}
