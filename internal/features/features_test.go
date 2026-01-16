package features

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

// TestDefaultFeatures verifies the default feature configuration
func TestDefaultFeatures(t *testing.T) {
	features := DefaultFeatures()

	if !features.GracefulShutdown {
		t.Error("GracefulShutdown should be enabled by default")
	}
	if !features.CertificateWatcher {
		t.Error("CertificateWatcher should be enabled by default")
	}
	if !features.PeriodicCertCheck {
		t.Error("PeriodicCertCheck should be enabled by default")
	}
	if !features.Logging {
		t.Error("Logging should be enabled by default")
	}
	if features.MetricsCollection {
		t.Error("MetricsCollection should be disabled by default")
	}
	if features.ShutdownTimeout != 10 {
		t.Errorf("ShutdownTimeout should be 10, got %d", features.ShutdownTimeout)
	}
	if features.AgentShutdownTimeout != 5 {
		t.Errorf("AgentShutdownTimeout should be 5, got %d", features.AgentShutdownTimeout)
	}
}

// TestMinimalFeatures verifies the minimal feature configuration
func TestMinimalFeatures(t *testing.T) {
	features := MinimalFeatures()

	if features.GracefulShutdown {
		t.Error("GracefulShutdown should be disabled in minimal mode")
	}
	if features.CertificateWatcher {
		t.Error("CertificateWatcher should be disabled in minimal mode")
	}
	if !features.Logging {
		t.Error("Logging should remain enabled in minimal mode")
	}
}

// TestAllFeatures verifies all features are enabled
func TestAllFeatures(t *testing.T) {
	features := AllFeatures()

	if !features.GracefulShutdown {
		t.Error("GracefulShutdown should be enabled")
	}
	if !features.CertificateWatcher {
		t.Error("CertificateWatcher should be enabled")
	}
	if !features.MetricsCollection {
		t.Error("MetricsCollection should be enabled")
	}
	if !features.HealthCheck {
		t.Error("HealthCheck should be enabled")
	}
}

// TestNewConfigLoader creates a new config loader with defaults
func TestNewConfigLoader(t *testing.T) {
	loader := NewConfigLoader()

	if loader == nil {
		t.Fatal("NewConfigLoader should not return nil")
	}

	features := loader.Get()
	if !features.GracefulShutdown {
		t.Error("New loader should have default features")
	}
}

// TestLoadFromEnv loads features from environment variables
func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("TLS_AGENT_FEATURES_GRACEFUL_SHUTDOWN", "false")
	os.Setenv("TLS_AGENT_FEATURES_LOGGING", "true")
	os.Setenv("TLS_AGENT_FEATURES_SHUTDOWN_TIMEOUT", "15")
	os.Setenv("TLS_AGENT_FEATURES_CERT_WATCH_INTERVAL", "60")

	defer func() {
		os.Unsetenv("TLS_AGENT_FEATURES_GRACEFUL_SHUTDOWN")
		os.Unsetenv("TLS_AGENT_FEATURES_LOGGING")
		os.Unsetenv("TLS_AGENT_FEATURES_SHUTDOWN_TIMEOUT")
		os.Unsetenv("TLS_AGENT_FEATURES_CERT_WATCH_INTERVAL")
	}()

	loader := NewConfigLoader()
	err := loader.LoadFromEnv()

	if err != nil {
		t.Fatalf("LoadFromEnv should not return error: %v", err)
	}

	features := loader.Get()

	if features.GracefulShutdown {
		t.Error("GracefulShutdown should be false from env")
	}
	if !features.Logging {
		t.Error("Logging should be true from env")
	}
	if features.ShutdownTimeout != 15 {
		t.Errorf("ShutdownTimeout should be 15 from env, got %d", features.ShutdownTimeout)
	}
	if features.CertWatchInterval != 60 {
		t.Errorf("CertWatchInterval should be 60 from env, got %d", features.CertWatchInterval)
	}
}

// TestLoadFromJSON loads features from a JSON file
func TestLoadFromJSON(t *testing.T) {
	// Create a temporary JSON file
	jsonData := map[string]interface{}{
		"graceful_shutdown":      true,
		"certificate_watcher":    false,
		"logging":                true,
		"shutdown_timeout":       12,
		"agent_shutdown_timeout": 4,
		"cert_watch_interval":    45,
		"debounce_interval":      1500,
		"cert_expiry_warning":    10,
	}

	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "features_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	loader := NewConfigLoader()
	err = loader.LoadFromJSON(tmpFile.Name())

	if err != nil {
		t.Fatalf("LoadFromJSON should not return error: %v", err)
	}

	features := loader.Get()

	if !features.GracefulShutdown {
		t.Error("GracefulShutdown should be true from JSON")
	}
	if features.CertificateWatcher {
		t.Error("CertificateWatcher should be false from JSON")
	}
	if features.ShutdownTimeout != 12 {
		t.Errorf("ShutdownTimeout should be 12 from JSON, got %d", features.ShutdownTimeout)
	}
	if features.CertWatchInterval != 45 {
		t.Errorf("CertWatchInterval should be 45 from JSON, got %d", features.CertWatchInterval)
	}
}

// TestLoadFromYAML loads features from a YAML file
func TestLoadFromYAML(t *testing.T) {
	yamlContent := `
graceful_shutdown: true
certificate_watcher: false
periodic_cert_check: true
debounce_file_changes: false
logging: true
metrics_collection: false
health_check: false
shutdown_timeout: 20
agent_shutdown_timeout: 8
cert_watch_interval: 40
debounce_interval: 3000
cert_expiry_warning: 14
`

	tmpFile, err := os.CreateTemp("", "features_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(yamlContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	loader := NewConfigLoader()
	err = loader.LoadFromYAML(tmpFile.Name())

	if err != nil {
		t.Fatalf("LoadFromYAML should not return error: %v", err)
	}

	features := loader.Get()

	if !features.GracefulShutdown {
		t.Error("GracefulShutdown should be true from YAML")
	}
	if features.CertificateWatcher {
		t.Error("CertificateWatcher should be false from YAML")
	}
	if features.ShutdownTimeout != 20 {
		t.Errorf("ShutdownTimeout should be 20 from YAML, got %d", features.ShutdownTimeout)
	}
}

// TestConfigLoaderSet replaces entire configuration
func TestConfigLoaderSet(t *testing.T) {
	loader := NewConfigLoader()

	customFeatures := Features{
		GracefulShutdown:     false,
		CertificateWatcher:   false,
		Logging:              true,
		ShutdownTimeout:      30,
		AgentShutdownTimeout: 10,
	}

	loader.Set(customFeatures)
	features := loader.Get()

	if features.GracefulShutdown {
		t.Error("GracefulShutdown should be false after Set")
	}
	if features.ShutdownTimeout != 30 {
		t.Errorf("ShutdownTimeout should be 30 after Set, got %d", features.ShutdownTimeout)
	}
}

// TestConfigLoaderUpdate modifies individual features
func TestConfigLoaderUpdate(t *testing.T) {
	loader := NewConfigLoader()

	loader.Update("graceful_shutdown", false)
	loader.Update("shutdown_timeout", 25)
	loader.Update("logging", false)

	features := loader.Get()

	if features.GracefulShutdown {
		t.Error("GracefulShutdown should be false after Update")
	}
	if features.ShutdownTimeout != 25 {
		t.Errorf("ShutdownTimeout should be 25 after Update, got %d", features.ShutdownTimeout)
	}
	if features.Logging {
		t.Error("Logging should be false after Update")
	}
}

// TestLoadFromNonexistentFile handles missing files gracefully
func TestLoadFromNonexistentFile(t *testing.T) {
	loader := NewConfigLoader()

	err := loader.LoadFromJSON("/nonexistent/path/features.json")
	if err == nil {
		t.Error("LoadFromJSON should return error for nonexistent file")
	}

	err = loader.LoadFromYAML("/nonexistent/path/features.yaml")
	if err == nil {
		t.Error("LoadFromYAML should return error for nonexistent file")
	}
}

// TestFeatureTimeouts ensures timeout values work correctly
func TestFeatureTimeouts(t *testing.T) {
	loader := NewConfigLoader()

	loader.Update("shutdown_timeout", 5)
	loader.Update("agent_shutdown_timeout", 2)

	features := loader.Get()

	// Verify timeout values can be used with time.Duration
	shutdownCtx := time.Duration(features.ShutdownTimeout) * time.Second
	agentCtx := time.Duration(features.AgentShutdownTimeout) * time.Second

	if shutdownCtx != 5*time.Second {
		t.Errorf("ShutdownTimeout conversion failed, got %v", shutdownCtx)
	}
	if agentCtx != 2*time.Second {
		t.Errorf("AgentShutdownTimeout conversion failed, got %v", agentCtx)
	}
}

// TestEnvironmentVariablePriority verifies env vars override config files
func TestEnvironmentVariablePriority(t *testing.T) {
	// Set environment variable to override
	os.Setenv("TLS_AGENT_FEATURES_SHUTDOWN_TIMEOUT", "99")
	defer os.Unsetenv("TLS_AGENT_FEATURES_SHUTDOWN_TIMEOUT")

	jsonData := map[string]interface{}{
		"graceful_shutdown": true,
		"shutdown_timeout":  10, // This should be overridden
	}

	data, _ := json.MarshalIndent(jsonData, "", "  ")

	tmpFile, _ := os.CreateTemp("", "features_*.json")
	defer os.Remove(tmpFile.Name())

	tmpFile.Write(data)
	tmpFile.Close()

	loader := NewConfigLoader()
	loader.LoadFromJSON(tmpFile.Name())
	loader.LoadFromEnv() // This should override

	features := loader.Get()

	if features.ShutdownTimeout != 99 {
		t.Errorf("Environment variable should override JSON config, got %d", features.ShutdownTimeout)
	}
}
