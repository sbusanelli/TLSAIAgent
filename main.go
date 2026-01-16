package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tls-agent/internal/agent"
	"tls-agent/internal/tlsstore"
)

func main() {
	cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatal(err)
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
		Addr:      ":8443",
		TLSConfig: tlsCfg,
	}

	// Channel for graceful shutdown
	shutdownDone := make(chan struct{})

	// Handle signals in a goroutine
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		log.Println("Initiating graceful shutdown...")

		// Signal the agent to stop
		close(agentStopChan)

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the HTTP server
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		log.Println("Server shutdown complete")
		close(shutdownDone)
	}()

	log.Println("TLS Agent server running on :8443")
	log.Println("Press Ctrl+C to gracefully shutdown")

	if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
		log.Printf("Server error: %v", err)
	}

	// Wait for shutdown to complete
	<-shutdownDone

	// Wait for agent to stop (with timeout)
	log.Println("Waiting for certificate watcher agent to stop...")
	agentStopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-agentDone:
		log.Println("Agent stopped gracefully")
	case <-agentStopCtx.Done():
		log.Println("Warning: Agent stop timeout (continuing anyway)")
	}

	log.Println("TLS Agent shutdown complete")
}
