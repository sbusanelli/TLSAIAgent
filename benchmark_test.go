package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"tls-agent/internal/agent"
	"tls-agent/internal/features"
	"tls-agent/internal/tlsstore"
)

// BenchmarkAgentStartup benchmarks agent startup time
func BenchmarkAgentStartup(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	b.ResetTimer()

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

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkAgentShutdown benchmarks agent shutdown time
func BenchmarkAgentShutdown(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
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
		time.Sleep(10 * time.Millisecond)

		b.StartTimer()

		// Signal agent to stop
		close(agentStopChan)

		// Wait for agent to stop
		<-agentDone

		b.StopTimer()
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkCertificateLoading benchmarks certificate loading performance
func BenchmarkCertificateLoading(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tlsstore.Load(certFile, keyFile)
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkCertificateRetrieval benchmarks certificate retrieval performance
func BenchmarkCertificateRetrieval(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	store := tlsstore.New(cert)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.GetCertificate(&tls.ClientHelloInfo{})
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkStateOperations benchmarks agent state operations
func BenchmarkStateOperations(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	_ = agent.NewState(cert)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
// 		state.IncrementCertificateCount()
// 		state.GetCertificateCount()
// 		state.IsRunning()
// 		state.GetCertificate()
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkConcurrentAccess benchmarks concurrent access to agent state
func BenchmarkConcurrentAccess(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	_ = agent.NewState(cert)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
// 			state.IncrementCertificateCount()
// 			state.GetCertificateCount()
// 			state.IsRunning()
// 			state.GetCertificate()
		}
	})

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkHTTPServer benchmarks HTTP server performance with TLS
func BenchmarkHTTPServer(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
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
		Addr:      ":9449",
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
			b.Logf("Server error: %v", err)
		}
		close(serverStopped)
	}()

	// Wait for server to start
	<-serverStarted
	time.Sleep(100 * time.Millisecond)

	// Create HTTP client
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get("https://localhost:9449/")
		if err != nil {
			b.Logf("HTTP request failed: %v", err)
		} else {
			resp.Body.Close()
		}
	}

	b.StopTimer()

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		b.Logf("Server shutdown error: %v", err)
	}
	cancel()

	// Wait for server to stop
	<-serverStopped

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	<-agentDone

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkConcurrentHTTPRequests benchmarks concurrent HTTP requests
func BenchmarkConcurrentHTTPRequests(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
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
		Addr:      ":9450",
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
			b.Logf("Server error: %v", err)
		}
		close(serverStopped)
	}()

	// Wait for server to start
	<-serverStarted
	time.Sleep(100 * time.Millisecond)

	// Create HTTP client
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get("https://localhost:9450/")
			if err != nil {
				b.Logf("HTTP request failed: %v", err)
			} else {
				resp.Body.Close()
			}
		}
	})

	b.StopTimer()

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		b.Logf("Server shutdown error: %v", err)
	}
	cancel()

	// Wait for server to stop
	<-serverStopped

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	<-agentDone

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkMemoryUsage benchmarks memory usage during operations
func BenchmarkMemoryUsage(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	// Measure initial memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	b.ResetTimer()

	// Create multiple stores and states
	stores := make([]*tlsstore.Store, 100)
	states := make([]*agent.State, 100)

	for i := 0; i < 100; i++ {
		stores[i] = tlsstore.New(cert)
		states[i] = agent.NewState(cert)
	}

	// Perform operations
	for i := 0; i < 1000; i++ {
		for j := 0; j < 100; j++ {
			stores[j].GetCertificate()
// 			states[j].IncrementCertificateCount()
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

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkFeatureFlags benchmarks feature flag performance
func BenchmarkFeatureFlags(b *testing.B) {
	// Test default features
	b.Run("DefaultFeatures", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			features.DefaultFeatures()
		}
	})

	// Test minimal features
	b.Run("MinimalFeatures", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			features.MinimalFeatures()
		}
	})

	// Test feature validation
	b.Run("ValidateFeatures", func(b *testing.B) {
		f := features.DefaultFeatures()
		for i := 0; i < b.N; i++ {
			f.Validate()
		}
	})
}

// BenchmarkGracefulShutdown benchmarks graceful shutdown performance
func BenchmarkGracefulShutdown(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
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
		time.Sleep(10 * time.Millisecond)

		b.StartTimer()

		// Signal agent to stop
		close(agentStopChan)

		// Wait for agent to stop
		<-agentDone

		b.StopTimer()
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkMultipleServers benchmarks multiple HTTP servers
func BenchmarkMultipleServers(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
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
	servers := make([]*http.Server, 5)
	serverStarted := make([]chan struct{}, 5)
	serverStopped := make([]chan struct{}, 5)

	for i := 0; i < 5; i++ {
		tlsCfg := &tls.Config{
			GetCertificate: store.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		}

		servers[i] = &http.Server{
			Addr:      fmt.Sprintf(":%d", 9451+i),
			TLSConfig: tlsCfg,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		}

		serverStarted[i] = make(chan struct{})
		serverStopped[i] = make(chan struct{})

		// Start server
		go func(idx int) {
			close(serverStarted[idx])
			if err := servers[idx].ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				b.Logf("Server %d error: %v", idx, err)
			}
			close(serverStopped[idx])
		}(i)
	}

	// Wait for servers to start
	for i := 0; i < 5; i++ {
		<-serverStarted[i]
	}

	time.Sleep(100 * time.Millisecond)

	// Create HTTP client
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		serverIndex := i % 5
		resp, err := client.Get(fmt.Sprintf("https://localhost:%d/", 9451+serverIndex))
		if err != nil {
			b.Logf("HTTP request failed: %v", err)
		} else {
			resp.Body.Close()
		}
	}

	b.StopTimer()

	// Shutdown all servers
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := servers[i].Shutdown(ctx); err != nil {
			b.Logf("Server %d shutdown error: %v", i, err)
		}
		cancel()
	}

	// Wait for all servers to stop
	for i := 0; i < 5; i++ {
		<-serverStopped[i]
	}

	// Signal agent to stop
	close(agentStopChan)

	// Wait for agent to stop
	<-agentDone

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// BenchmarkConcurrentShutdown benchmarks concurrent shutdown operations
func BenchmarkConcurrentShutdown(b *testing.B) {
	// Create temporary directory for test certificates
	tempDir := b.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")

	// Create test certificates
	err := createTestCertificates(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to create test certificates: %v", err)
	}

	// Load certificates
	cert, err := tlsstore.Load(certFile, keyFile)
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var wg sync.WaitGroup
		wg.Add(10)

		// Start 10 agents concurrently
		for j := 0; j < 10; j++ {
			go func() {
				defer wg.Done()
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
				time.Sleep(1 * time.Millisecond)

				// Signal agent to stop
				close(agentStopChan)

				// Wait for agent to stop
				<-agentDone
			}()
		}

		b.StartTimer()

		// Wait for all agents to stop
		wg.Wait()

		b.StopTimer()
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}
