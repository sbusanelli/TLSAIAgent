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
	"tls-agent/internal/features"
	"tls-agent/internal/tlsstore"
)

func main() {
	// Load feature configuration
	featureLoader := features.NewConfigLoader()

	// Try to load from config file if specified
	if configPath := os.Getenv("FEATURES_CONFIG_PATH"); configPath != "" {
		if err := featureLoader.LoadFromYAML(configPath); err != nil {
			if err := featureLoader.LoadFromJSON(configPath); err != nil {
				log.Printf("Warning: Could not load features config from %s: %v\n", configPath, err)
			}
		}
	}

	// Override with environment variables (takes precedence)
	if err := featureLoader.LoadFromEnv(); err != nil {
		log.Printf("Warning: Could not load features from environment: %v\n", err)
	}

	featureConfig := featureLoader.Get()
	featureLoader.LogFeatures()

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

	// Only start the certificate watcher agent if feature is enabled
	if featureConfig.CertificateWatcher {
		go func() {
			agent.Run(store, state, agentStopChan)
			close(agentDone)
		}()
	} else {
		close(agentDone) // Mark as already done if feature is disabled
		if featureConfig.Logging {
			log.Println("Certificate watcher agent disabled")
		}
	}

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsCfg,
	}

	// Channel for graceful shutdown
	shutdownDone := make(chan struct{})

	if featureConfig.GracefulShutdown {
		// Handle signals in a goroutine for graceful shutdown
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			sig := <-sigChan
			if featureConfig.Logging {
				log.Printf("Received signal: %v", sig)
				log.Println("Initiating graceful shutdown...")
			}

			// Signal the agent to stop
			if featureConfig.CertificateWatcher {
				close(agentStopChan)
			}

			// Create context with timeout for shutdown
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(featureConfig.ShutdownTimeout)*time.Second)
			defer cancel()

			// Shutdown the HTTP server
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("Server shutdown error: %v", err)
			}

			if featureConfig.Logging {
				log.Println("Server shutdown complete")
			}
			close(shutdownDone)
		}()
	} else {
		close(shutdownDone) // Mark as done if graceful shutdown is disabled
		if featureConfig.Logging {
			log.Println("Graceful shutdown feature disabled")
		}
	}

	if featureConfig.Logging {
		log.Println("TLS Agent server running on :8443")
		log.Println("Press Ctrl+C to gracefully shutdown")
	}

	if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
		log.Printf("Server error: %v", err)
	}

	// Wait for shutdown to complete
	<-shutdownDone

	// Wait for agent to stop (with timeout) if watcher is enabled
	if featureConfig.CertificateWatcher {
		if featureConfig.Logging {
			log.Println("Waiting for certificate watcher agent to stop...")
		}
		agentStopCtx, cancel := context.WithTimeout(context.Background(), time.Duration(featureConfig.AgentShutdownTimeout)*time.Second)
		defer cancel()

		select {
		case <-agentDone:
			if featureConfig.Logging {
				log.Println("Agent stopped gracefully")
			}
		case <-agentStopCtx.Done():
			log.Println("Warning: Agent stop timeout (continuing anyway)")
		}
	}

	log.Println("TLS Agent shutdown complete")
}
