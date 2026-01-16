package features

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Features represents all configurable features in the TLS Agent
type Features struct {
	// GracefulShutdown enables graceful shutdown with signal handling
	GracefulShutdown bool `json:"graceful_shutdown" yaml:"graceful_shutdown"`

	// CertificateWatcher enables the file-based certificate watcher agent
	CertificateWatcher bool `json:"certificate_watcher" yaml:"certificate_watcher"`

	// PeriodicCertCheck enables periodic certificate expiry checking
	PeriodicCertCheck bool `json:"periodic_cert_check" yaml:"periodic_cert_check"`

	// DebounceFileChanges enables debouncing of rapid certificate file changes
	DebounceFileChanges bool `json:"debounce_file_changes" yaml:"debounce_file_changes"`

	// Logging enables detailed logging throughout the application
	Logging bool `json:"logging" yaml:"logging"`

	// MetricsCollection enables collection of performance metrics
	MetricsCollection bool `json:"metrics_collection" yaml:"metrics_collection"`

	// HealthCheck enables a health check endpoint (future feature)
	HealthCheck bool `json:"health_check" yaml:"health_check"`

	// ShutdownTimeout is the timeout duration for graceful shutdown in seconds
	ShutdownTimeout int `json:"shutdown_timeout" yaml:"shutdown_timeout"`

	// AgentShutdownTimeout is the timeout for agent shutdown in seconds
	AgentShutdownTimeout int `json:"agent_shutdown_timeout" yaml:"agent_shutdown_timeout"`

	// CertWatchInterval is the periodic check interval in seconds
	CertWatchInterval int `json:"cert_watch_interval" yaml:"cert_watch_interval"`

	// DebounceInterval is the debounce interval in milliseconds
	DebounceInterval int `json:"debounce_interval" yaml:"debounce_interval"`

	// CertExpiryWarning is the days before expiry to warn about certificate
	CertExpiryWarning int `json:"cert_expiry_warning" yaml:"cert_expiry_warning"`
}

// DefaultFeatures returns the default feature configuration with all features enabled
func DefaultFeatures() Features {
	return Features{
		GracefulShutdown:     true,
		CertificateWatcher:   true,
		PeriodicCertCheck:    true,
		DebounceFileChanges:  true,
		Logging:              true,
		MetricsCollection:    false, // Disabled by default (future feature)
		HealthCheck:          false, // Disabled by default (future feature)
		ShutdownTimeout:      10,
		AgentShutdownTimeout: 5,
		CertWatchInterval:    30,
		DebounceInterval:     2000, // 2 seconds in milliseconds
		CertExpiryWarning:    7,    // 7 days
	}
}

// MinimalFeatures returns a minimal configuration with only core features enabled
func MinimalFeatures() Features {
	return Features{
		GracefulShutdown:     false,
		CertificateWatcher:   false,
		PeriodicCertCheck:    false,
		DebounceFileChanges:  false,
		Logging:              true,
		MetricsCollection:    false,
		HealthCheck:          false,
		ShutdownTimeout:      5,
		AgentShutdownTimeout: 2,
		CertWatchInterval:    60,
		DebounceInterval:     1000,
		CertExpiryWarning:    14,
	}
}

// AllFeatures returns a configuration with all features enabled
func AllFeatures() Features {
	return Features{
		GracefulShutdown:     true,
		CertificateWatcher:   true,
		PeriodicCertCheck:    true,
		DebounceFileChanges:  true,
		Logging:              true,
		MetricsCollection:    true,
		HealthCheck:          true,
		ShutdownTimeout:      10,
		AgentShutdownTimeout: 5,
		CertWatchInterval:    30,
		DebounceInterval:     2000,
		CertExpiryWarning:    7,
	}
}

// ConfigLoader provides methods to load feature configurations from various sources
type ConfigLoader struct {
	features Features
}

// NewConfigLoader creates a new configuration loader with default features
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		features: DefaultFeatures(),
	}
}

// LoadFromEnv loads feature flags from environment variables
// Environment variable format: TLS_AGENT_FEATURES_<FEATURE_NAME>=true/false
// Example: TLS_AGENT_FEATURES_GRACEFUL_SHUTDOWN=true
func (cl *ConfigLoader) LoadFromEnv() error {
	// Load boolean features
	cl.loadBoolEnv("GRACEFUL_SHUTDOWN", &cl.features.GracefulShutdown)
	cl.loadBoolEnv("CERTIFICATE_WATCHER", &cl.features.CertificateWatcher)
	cl.loadBoolEnv("PERIODIC_CERT_CHECK", &cl.features.PeriodicCertCheck)
	cl.loadBoolEnv("DEBOUNCE_FILE_CHANGES", &cl.features.DebounceFileChanges)
	cl.loadBoolEnv("LOGGING", &cl.features.Logging)
	cl.loadBoolEnv("METRICS_COLLECTION", &cl.features.MetricsCollection)
	cl.loadBoolEnv("HEALTH_CHECK", &cl.features.HealthCheck)

	// Load integer features
	cl.loadIntEnv("SHUTDOWN_TIMEOUT", &cl.features.ShutdownTimeout)
	cl.loadIntEnv("AGENT_SHUTDOWN_TIMEOUT", &cl.features.AgentShutdownTimeout)
	cl.loadIntEnv("CERT_WATCH_INTERVAL", &cl.features.CertWatchInterval)
	cl.loadIntEnv("DEBOUNCE_INTERVAL", &cl.features.DebounceInterval)
	cl.loadIntEnv("CERT_EXPIRY_WARNING", &cl.features.CertExpiryWarning)

	return nil
}

// LoadFromYAML loads feature flags from a YAML configuration file
func (cl *ConfigLoader) LoadFromYAML(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &cl.features)
	if err != nil {
		return err
	}

	if cl.features.Logging {
		log.Printf("Features loaded from YAML file: %s\n", filePath)
	}

	return nil
}

// LoadFromJSON loads feature flags from a JSON configuration file
func (cl *ConfigLoader) LoadFromJSON(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &cl.features)
	if err != nil {
		return err
	}

	if cl.features.Logging {
		log.Printf("Features loaded from JSON file: %s\n", filePath)
	}

	return nil
}

// Get returns the current feature configuration
func (cl *ConfigLoader) Get() Features {
	return cl.features
}

// Set replaces the entire feature configuration
func (cl *ConfigLoader) Set(features Features) {
	cl.features = features
}

// Update modifies a specific feature flag
func (cl *ConfigLoader) Update(featureName string, value interface{}) {
	switch strings.ToLower(featureName) {
	case "graceful_shutdown":
		if b, ok := value.(bool); ok {
			cl.features.GracefulShutdown = b
		}
	case "certificate_watcher":
		if b, ok := value.(bool); ok {
			cl.features.CertificateWatcher = b
		}
	case "periodic_cert_check":
		if b, ok := value.(bool); ok {
			cl.features.PeriodicCertCheck = b
		}
	case "debounce_file_changes":
		if b, ok := value.(bool); ok {
			cl.features.DebounceFileChanges = b
		}
	case "logging":
		if b, ok := value.(bool); ok {
			cl.features.Logging = b
		}
	case "metrics_collection":
		if b, ok := value.(bool); ok {
			cl.features.MetricsCollection = b
		}
	case "health_check":
		if b, ok := value.(bool); ok {
			cl.features.HealthCheck = b
		}
	case "shutdown_timeout":
		if i, ok := value.(int); ok {
			cl.features.ShutdownTimeout = i
		}
	case "agent_shutdown_timeout":
		if i, ok := value.(int); ok {
			cl.features.AgentShutdownTimeout = i
		}
	}
}

// LogFeatures logs all enabled features
func (cl *ConfigLoader) LogFeatures() {
	if !cl.features.Logging {
		return
	}

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("Feature Configuration:")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("  Graceful Shutdown:     %v\n", cl.features.GracefulShutdown)
	log.Printf("  Certificate Watcher:   %v\n", cl.features.CertificateWatcher)
	log.Printf("  Periodic Cert Check:   %v\n", cl.features.PeriodicCertCheck)
	log.Printf("  Debounce File Changes: %v\n", cl.features.DebounceFileChanges)
	log.Printf("  Logging:               %v\n", cl.features.Logging)
	log.Printf("  Metrics Collection:    %v\n", cl.features.MetricsCollection)
	log.Printf("  Health Check:          %v\n", cl.features.HealthCheck)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("  Shutdown Timeout:      %d seconds\n", cl.features.ShutdownTimeout)
	log.Printf("  Agent Shutdown Timeout: %d seconds\n", cl.features.AgentShutdownTimeout)
	log.Printf("  Cert Watch Interval:   %d seconds\n", cl.features.CertWatchInterval)
	log.Printf("  Debounce Interval:     %d ms\n", cl.features.DebounceInterval)
	log.Printf("  Cert Expiry Warning:   %d days\n", cl.features.CertExpiryWarning)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// Helper functions

func (cl *ConfigLoader) loadBoolEnv(envName string, target *bool) {
	fullEnvName := "TLS_AGENT_FEATURES_" + envName
	if val, exists := os.LookupEnv(fullEnvName); exists {
		if parsedVal, err := strconv.ParseBool(val); err == nil {
			*target = parsedVal
		}
	}
}

func (cl *ConfigLoader) loadIntEnv(envName string, target *int) {
	fullEnvName := "TLS_AGENT_FEATURES_" + envName
	if val, exists := os.LookupEnv(fullEnvName); exists {
		if parsedVal, err := strconv.Atoi(val); err == nil {
			*target = parsedVal
		}
	}
}
