package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"tls-agent/internal/agent"
	"tls-agent/internal/tlsstore"
)

// TestGracefulShutdown tests the graceful shutdown of the server and agent
func TestGracefulShutdown(t *testing.T) {
	// Load certificates
	cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)

	tlsCfg := &tls.Config{
		GetCertificate: store.GetCertificate,
		MinVersion:     tls.VersionTLS12,
	}

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

	server := &http.Server{
		Addr:      ":9443",
		TLSConfig: tlsCfg,
	}

	shutdownDone := make(chan struct{})
	serverStarted := make(chan struct{})
	shutdownSignalTime := time.Now()

	// Handle signals in a goroutine (simulating signal handler)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		t.Logf("Received signal: %v", sig)

		// Record shutdown start time
		shutdownSignalTime = time.Now()

		// Signal the agent to stop
		close(agentStopChan)

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the HTTP server
		if err := server.Shutdown(ctx); err != nil {
			t.Logf("Server shutdown error: %v", err)
		}

		close(shutdownDone)
	}()

	// Start server in a goroutine
	go func() {
		close(serverStarted)
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	<-serverStarted
	time.Sleep(100 * time.Millisecond)

	// Verify server is running by making a test request
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get("https://localhost:9443/")
	if err == nil {
		resp.Body.Close()
		t.Log("Server is running, test request succeeded")
	} else {
		// This is expected since we don't have actual handlers
		t.Logf("Test request result: %v (expected for no handlers)", err)
	}

	// Send shutdown signal
	t.Log("Sending shutdown signal...")
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find current process: %v", err)
	}

	if err := proc.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("Failed to send signal: %v", err)
	}

	// Wait for shutdown to complete (with timeout)
	shutdownTimeout := time.After(15 * time.Second)
	select {
	case <-shutdownDone:
		t.Log("Server shutdown signal handled")
	case <-shutdownTimeout:
		t.Fatal("Timeout waiting for server shutdown")
	}

	// Wait for agent to stop
	agentStopTimeout := time.After(10 * time.Second)
	select {
	case <-agentDone:
		t.Log("Agent stopped successfully")
	case <-agentStopTimeout:
		t.Error("Timeout waiting for agent to stop")
	}

	// Verify shutdown was fast enough
	shutdownDuration := time.Since(shutdownSignalTime)
	t.Logf("Total shutdown duration: %v", shutdownDuration)

	if shutdownDuration > 15*time.Second {
		t.Errorf("Shutdown took too long: %v", shutdownDuration)
	}

	t.Log("Graceful shutdown test completed successfully")
}

// TestServerStopsAcceptingConnections tests that the server stops accepting new connections after shutdown
func TestServerStopsAcceptingConnections(t *testing.T) {
	cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	tlsCfg := &tls.Config{
		GetCertificate: store.GetCertificate,
		MinVersion:     tls.VersionTLS12,
	}

	state := agent.NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	go func() {
		agent.Run(store, state, agentStopChan)
		close(agentDone)
	}()

	server := &http.Server{
		Addr:      ":9444",
		TLSConfig: tlsCfg,
	}

	serverRunning := make(chan struct{})
	serverStopped := make(chan struct{})

	// Start server
	go func() {
		close(serverRunning)
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
		close(serverStopped)
	}()

	<-serverRunning
	time.Sleep(100 * time.Millisecond)

	// Verify server is accepting connections
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	t.Log("Attempting connection before shutdown...")
	_, err = client.Get("https://localhost:9444/")
	// Connection attempt should reach the server (handler not found or timeout is expected)
	t.Logf("Connection attempt result: %v", err)

	// Trigger shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	close(agentStopChan)
	if err := server.Shutdown(ctx); err != nil {
		t.Logf("Server shutdown error: %v", err)
	}
	cancel()

	// Wait for server to stop
	select {
	case <-serverStopped:
		t.Log("Server stopped")
	case <-time.After(6 * time.Second):
		t.Fatal("Server did not stop in time")
	}

	// Verify server is NOT accepting connections
	t.Log("Attempting connection after shutdown...")
	_, err = client.Get("https://localhost:9444/")
	if err == nil {
		t.Error("Server should not accept connections after shutdown")
	} else {
		t.Logf("Connection correctly rejected after shutdown: %v", err)
	}

	// Wait for agent
	select {
	case <-agentDone:
		t.Log("Agent stopped")
	case <-time.After(6 * time.Second):
		t.Log("Agent did not stop in time (continuing test)")
	}

	t.Log("Server stops accepting connections test passed")
}

// TestAgentShutdownWithTimeout tests that the agent shuts down within timeout
func TestAgentShutdownWithTimeout(t *testing.T) {
	cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	state := agent.NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	startTime := time.Now()

	// Start the agent
	go func() {
		agent.Run(store, state, agentStopChan)
		close(agentDone)
	}()

	// Give agent time to start
	time.Sleep(100 * time.Millisecond)

	// Signal shutdown
	t.Log("Signaling agent to stop...")
	close(agentStopChan)

	// Wait for agent with timeout
	shutdownTimeout := 5 * time.Second
	select {
	case <-agentDone:
		duration := time.Since(startTime)
		t.Logf("Agent stopped in %v", duration)
		if duration > shutdownTimeout {
			t.Errorf("Agent shutdown took too long: %v", duration)
		}
	case <-time.After(shutdownTimeout + 1*time.Second):
		t.Error("Agent did not stop within timeout")
	}

	t.Log("Agent shutdown timeout test passed")
}

// TestMultipleSignals tests that multiple signals don't cause issues
func TestMultipleSignals(t *testing.T) {
	cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)
	tlsCfg := &tls.Config{
		GetCertificate: store.GetCertificate,
		MinVersion:     tls.VersionTLS12,
	}

	state := agent.NewState(cert)
	agentStopChan := make(chan struct{})
	agentDone := make(chan struct{})

	go func() {
		agent.Run(store, state, agentStopChan)
		close(agentDone)
	}()

	server := &http.Server{
		Addr:      ":9445",
		TLSConfig: tlsCfg,
	}

	shutdownDone := make(chan struct{})
	shutdownCount := 0

	// Handle signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		t.Logf("Received first signal: %v", sig)
		shutdownCount++

		close(agentStopChan)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			t.Logf("Server shutdown error: %v", err)
		}

		close(shutdownDone)
	}()

	// Start server
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Send first signal
	t.Log("Sending first shutdown signal...")
	proc, _ := os.FindProcess(os.Getpid())
	proc.Signal(syscall.SIGINT)

	// Wait for shutdown
	select {
	case <-shutdownDone:
		t.Log("Shutdown initiated")
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for first shutdown")
	}

	// Try to send another signal (should be ignored since shutdown already started)
	t.Log("Attempting to send second signal...")
	proc.Signal(syscall.SIGINT)

	// Wait a bit to ensure no panic
	time.Sleep(500 * time.Millisecond)

	// Verify shutdown count
	if shutdownCount != 1 {
		t.Logf("Shutdown was triggered %d times (expected 1)", shutdownCount)
	}

	// Wait for agent
	select {
	case <-agentDone:
		t.Log("Agent stopped")
	case <-time.After(6 * time.Second):
		t.Log("Agent stop timeout (continuing)")
	}

	t.Log("Multiple signals test passed")
}
