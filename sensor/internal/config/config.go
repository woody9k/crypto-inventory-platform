package config

import (
	"os"
	"strconv"
	"time"
)

// Config represents sensor configuration
type Config struct {
	// Sensor identity
	SensorID    string `json:"sensor_id"`
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Platform    string `json:"platform"`
	Version     string `json:"version"`
	Profile     string `json:"profile"`

	// Control plane connection
	ControlPlaneURL string `json:"control_plane_url"`
	RegistrationKey string `json:"registration_key"`

	// Reporting configuration
	ReportingInterval time.Duration `json:"reporting_interval"`
	BatchSize         int           `json:"batch_size"`

	// Storage configuration
	Storage StorageConfig `json:"storage"`

	// Capture configuration
	Capture CaptureConfig `json:"capture"`

	// Network configuration
	Network NetworkConfig `json:"network"`

	// Security configuration
	Security SecurityConfig `json:"security"`

	// Features
	Features map[string]bool `json:"features"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	MaxStorageSize int64  `json:"max_storage_size"` // bytes
	RotationSize   int64  `json:"rotation_size"`    // bytes
	RetentionDays  int    `json:"retention_days"`
	DataPath       string `json:"data_path"`
	EncryptionKey  string `json:"encryption_key"`
}

// CaptureConfig represents packet capture configuration
type CaptureConfig struct {
	Interfaces       []string `json:"interfaces"`
	ActiveProbing    bool     `json:"active_probing"`
	NetworkDiscovery bool     `json:"network_discovery"`
	MaxConnections   int      `json:"max_connections"`
	TimeoutSeconds   int      `json:"timeout_seconds"`
	BufferSize       int      `json:"buffer_size"`
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	Interfaces []string `json:"interfaces"`
	VLANs      []string `json:"vlans"`
	Gateways   []string `json:"gateways"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	ClientCert   string `json:"client_cert"`
	ClientKey    string `json:"client_key"`
	ServerCACert string `json:"server_ca_cert"`
	UseTLS       bool   `json:"use_tls"`
}

// Load loads configuration from environment variables and defaults
func Load() *Config {
	cfg := &Config{
		SensorID:          getEnv("SENSOR_ID", ""),
		TenantID:          getEnv("TENANT_ID", "default-tenant"),
		Name:              getEnv("SENSOR_NAME", "crypto-sensor"),
		Description:       getEnv("SENSOR_DESCRIPTION", "Crypto Inventory Network Sensor"),
		Platform:          getEnv("SENSOR_PLATFORM", "linux"),
		Version:           getEnv("SENSOR_VERSION", "1.0.0"),
		Profile:           getEnv("SENSOR_PROFILE", "datacenter_host"),
		ControlPlaneURL:   getEnv("CONTROL_PLANE_URL", "http://localhost:8085"),
		RegistrationKey:   getEnv("REGISTRATION_KEY", ""),
		ReportingInterval: getDurationEnv("REPORTING_INTERVAL", 30*time.Second),
		BatchSize:         getIntEnv("BATCH_SIZE", 100),
		Storage: StorageConfig{
			MaxStorageSize: getInt64Env("MAX_STORAGE_SIZE", 100*1024*1024), // 100MB
			RotationSize:   getInt64Env("ROTATION_SIZE", 10*1024*1024),     // 10MB
			RetentionDays:  getIntEnv("RETENTION_DAYS", 7),
			DataPath:       getEnv("DATA_PATH", "/var/lib/crypto-sensor"),
			EncryptionKey:  getEnv("ENCRYPTION_KEY", ""),
		},
		Capture: CaptureConfig{
			Interfaces:       getStringSliceEnv("INTERFACES", []string{"eth0"}),
			ActiveProbing:    getBoolEnv("ACTIVE_PROBING", true),
			NetworkDiscovery: getBoolEnv("NETWORK_DISCOVERY", true),
			MaxConnections:   getIntEnv("MAX_CONNECTIONS", 1000),
			TimeoutSeconds:   getIntEnv("TIMEOUT_SECONDS", 30),
			BufferSize:       getIntEnv("BUFFER_SIZE", 1024*1024), // 1MB
		},
		Network: NetworkConfig{
			Interfaces: getStringSliceEnv("NETWORK_INTERFACES", []string{"eth0"}),
			VLANs:      getStringSliceEnv("VLANS", []string{}),
			Gateways:   getStringSliceEnv("GATEWAYS", []string{}),
		},
		Security: SecurityConfig{
			ClientCert:   getEnv("CLIENT_CERT", ""),
			ClientKey:    getEnv("CLIENT_KEY", ""),
			ServerCACert: getEnv("SERVER_CA_CERT", ""),
			UseTLS:       getBoolEnv("USE_TLS", false),
		},
		Features: map[string]bool{
			"tls_analysis":         getBoolEnv("FEATURE_TLS_ANALYSIS", true),
			"ssh_analysis":         getBoolEnv("FEATURE_SSH_ANALYSIS", true),
			"certificate_analysis": getBoolEnv("FEATURE_CERTIFICATE_ANALYSIS", true),
			"active_probing":       getBoolEnv("FEATURE_ACTIVE_PROBING", true),
			"network_discovery":    getBoolEnv("FEATURE_NETWORK_DISCOVERY", true),
			"air_gapped_export":    getBoolEnv("FEATURE_AIR_GAPPED_EXPORT", false),
		},
	}

	return cfg
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated parsing
		// In a real implementation, you might want more sophisticated parsing
		return []string{value}
	}
	return defaultValue
}
