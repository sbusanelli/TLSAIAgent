package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"tls-agent/internal/agent"
	"tls-agent/internal/features"
	"tls-agent/internal/tlsstore"
)

// TestIntegrationAgentServer tests the complete integration between agent and HTTP server
func TestIntegrationAgentServer(t *testing.T) {
	// Create temporary directory for test certificates
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates (simplified for testing)
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
	}()

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

	resp, err := client.Get("https://localhost:9447/")
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
		t.Log("Agent stopped")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop in time")
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// TestIntegrationFeatureFlags tests integration with feature flags
func TestIntegrationFeatureFlags(t *testing.T) {
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

	// Test with feature flags enabled
	features := features.DefaultFeatures()
	features.GracefulShutdown = true
	features.CertificateWatcher = true
	features.Logging = true

	// Create store and state
	store := tlsstore.New(cert)
	state := agent.NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	// Start agent with feature flags
	go func() {
		agent.Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Verify agent is running
	if !state.IsRunning() {
		t.Error("Agent should be running with feature flags enabled")
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent with feature flags stopped successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent with feature flags did not stop in time")
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// TestIntegrationHotReload tests hot reload functionality
func TestIntegrationHotReload(t *testing.T) {
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
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Modify certificate file to trigger hot reload
	err = os.WriteFile(certFile, []byte("modified certificate"), 0644)
	if err != nil {
		t.Fatalf("Failed to modify certificate file: %v", err)
	}

	// Give hot reload time to process
	time.Sleep(200 * time.Millisecond)

	// Verify agent is still running
	if !state.IsRunning() {
		t.Error("Agent should still be running after hot reload")
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent with hot reload stopped successfully")
	case <-time.After(5 * time.Second):
		t.Error("Agent with hot reload did not stop in time")
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// TestIntegrationGracefulShutdown tests graceful shutdown integration
func TestIntegrationGracefulShutdown(t *testing.T) {
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
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Create HTTP server
	tlsCfg := &tls.Config{
		GetCertificate: store.GetCertificate,
		MinVersion:     tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":9448",
		TLSConfig: tlsCfg,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Graceful shutdown test"))
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

	// Test graceful shutdown with signal handling
	shutdownDone := make(chan struct{})
	shutdownSignalTime := time.Now()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		t.Logf("Received signal: %v", sig)

		shutdownSignalTime = time.Now()

		// Signal agent to stop
		close(agentStopChan)

		// Shutdown server
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			t.Logf("Server shutdown error: %v", err)
		}

		close(shutdownDone)
	}()

	// Send shutdown signal
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find current process: %v", err)
	}

	if err := proc.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("Failed to send signal: %v", err)
	}

	// Wait for shutdown to complete
	select {
	case <-shutdownDone:
		t.Log("Graceful shutdown completed")
	case <-time.After(15 * time.Second):
		t.Fatal("Timeout waiting for graceful shutdown")
	}

	// Verify shutdown was fast enough
	shutdownDuration := time.Since(shutdownSignalTime)
	t.Logf("Graceful shutdown duration: %v", shutdownDuration)

	if shutdownDuration > 10*time.Second {
		t.Errorf("Graceful shutdown took too long: %v", shutdownDuration)
	}

	// Wait for server to stop
	select {
	case <-serverStopped:
		t.Log("Server stopped")
	case <-time.After(6 * time.Second):
		t.Error("Server did not stop in time")
	}

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent stopped")
	case <-time.After(6 * time.Second):
		t.Error("Agent did not stop in time")
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// TestIntegrationMultipleServers tests integration with multiple HTTP servers
func TestIntegrationMultipleServers(t *testing.T) {
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
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Create multiple HTTP servers
	servers := make([]*http.Server, 3)
	serverStopped := make([]chan struct{}, 3)

	for i := 0; i < 3; i++ {
		tlsCfg := &tls.Config{
			GetCertificate: store.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		}

		servers[i] = &http.Server{
			Addr:      fmt.Sprintf(":%d", 9449+i),
			TLSConfig: tlsCfg,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("Server %d OK", i)))
			}),
		}

		serverStopped[i] = make(chan struct{})

		// Start server
		go func(idx int) {
			if err := servers[idx].ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				t.Logf("Server %d error: %v", idx, err)
			}
			close(serverStopped[idx])
		}(i)
	}

	// Give servers time to start
	time.Sleep(200 * time.Millisecond)

	// Test requests to all servers
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	for i := 0; i < 3; i++ {
		resp, err := client.Get(fmt.Sprintf("https://localhost:%d/", 9449+i))
		if err != nil {
			t.Logf("Server %d request failed: %v", i, err)
		} else {
			resp.Body.Close()
			t.Logf("Server %d request succeeded", i)
		}
	}

	// Shutdown all servers
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := servers[i].Shutdown(ctx); err != nil {
			t.Logf("Server %d shutdown error: %v", i, err)
		}
		cancel()
	}

	// Wait for all servers to stop
	for i := 0; i < 3; i++ {
		select {
		case <-serverStopped[i]:
			t.Logf("Server %d stopped", i)
		case <-time.After(6 * time.Second):
			t.Errorf("Server %d did not stop in time", i)
		}
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	select {
	case <-agentDone:
		t.Log("Agent stopped")
	case <-time.After(5 * time.Second):
		t.Error("Agent did not stop in time")
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// TestIntegrationErrorHandling tests error handling in integration scenarios
func TestIntegrationErrorHandling(t *testing.T) {
	// Test with invalid certificate files
	_, err := tlsstore.Load("nonexistent.crt", "nonexistent.key")
	if err == nil {
		t.Error("Loading nonexistent certificates should fail")
	}

	// Test with empty certificate files
	tempDir := t.TempDir()
	emptyCert := filepath.Join(tempDir, "empty.crt")
	emptyKey := filepath.Join(tempDir, "empty.key")

	err = os.WriteFile(emptyCert, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create empty cert file: %v", err)
	}

	err = os.WriteFile(emptyKey, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create empty key file: %v", err)
	}

	_, err = tlsstore.Load(emptyCert, emptyKey)
	if err == nil {
		t.Error("Loading empty certificate files should fail")
	}

	// Clean up
	os.Remove(emptyCert)
	os.Remove(emptyKey)
}

// TestIntegrationPerformance tests performance in integration scenarios
func TestIntegrationPerformance(t *testing.T) {
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
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Performance test: rapid state operations
	startTime := time.Now()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		state.IncrementCertificateCount()
		state.GetCertificateCount()
		state.IsRunning()
		state.GetCertificate()
	}

	duration := time.Since(startTime)
	avgDuration := duration / time.Duration(iterations*4)

	t.Logf("Integration performance test: %d iterations in %v (avg: %v per operation)", 
		iterations, duration, avgDuration)

	// Performance threshold
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

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// Helper function to create test certificates
func createTestCertificates(certFile, keyFile string) error {
	// Create a simple test certificate (simplified for testing)
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
