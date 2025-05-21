package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config contains configuration items for all services
type Config struct {
	Service  ServiceConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Search   SearchConfig
	NATS     NATSConfig
	Auth     AuthConfig
	Trace    TraceConfig
	HTTP     HTTPConfig
	GRPC     GRPCConfig
}

// ServiceConfig contains basic service information
type ServiceConfig struct {
	Name        string
	Environment string
	LogLevel    string
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// SearchConfig contains search engine configuration
type SearchConfig struct {
	URL       string
	APIKey    string
	IndexName string
}

// NATSConfig contains NATS configuration
type NATSConfig struct {
	URL string
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSecret     string
	TokenDuration int // minutes
}

// TraceConfig contains distributed tracing configuration
type TraceConfig struct {
	Enabled bool
	URL     string
}

// HTTPConfig contains HTTP server configuration
type HTTPConfig struct {
	Port    int
	Timeout int // seconds
}

// GRPCConfig contains gRPC server configuration
type GRPCConfig struct {
	Port int
}

// DSN returns PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// RedisAddr returns Redis address
func (c *RedisConfig) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Load loads configuration from file and environment variables
func Load(serviceName, configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v, serviceName)

	// Read config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.AddConfigPath("./configs")
		v.AddConfigPath("../configs")
		v.AddConfigPath("../../configs")
		v.SetConfigName(serviceName)
	}
	
	// Support environment variable override
	v.SetEnvPrefix("GOSHOP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// If config file not found, just warn, not error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// Set default configuration
func setDefaults(v *viper.Viper, serviceName string) {
	// Service configuration
	v.SetDefault("service.name", serviceName)
	v.SetDefault("service.environment", "development")
	v.SetDefault("service.logLevel", "info")

	// Database configuration
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "goshop")
	v.SetDefault("database.password", "goshop")
	v.SetDefault("database.dbname", fmt.Sprintf("goshop_%s", serviceName))
	v.SetDefault("database.sslmode", "disable")

	// Redis configuration
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	// Search engine configuration
	v.SetDefault("search.url", "http://localhost:7700")
	v.SetDefault("search.apiKey", "masterKey")
	v.SetDefault("search.indexName", fmt.Sprintf("%s_index", serviceName))

	// NATS configuration
	v.SetDefault("nats.url", "nats://localhost:4222")

	// Authentication configuration
	v.SetDefault("auth.jwtSecret", "change-me-in-production")
	v.SetDefault("auth.tokenDuration", 60) // 60 minutes

	// Tracing configuration
	v.SetDefault("trace.enabled", true)
	v.SetDefault("trace.url", "http://localhost:14268/api/traces")

	// HTTP configuration
	v.SetDefault("http.port", getDefaultHTTPPort(serviceName))
	v.SetDefault("http.timeout", 30) // 30 seconds

	// gRPC configuration
	v.SetDefault("grpc.port", getDefaultGRPCPort(serviceName))
}

// Assign unique default port for each service
func getDefaultHTTPPort(serviceName string) int {
	ports := map[string]int{
		"user":      8001,
		"product":   8002,
		"inventory": 8003,
		"order":     8004,
		"payment":   8005,
		"marketing": 8006,
		"cms":       8007,
		"shipping":  8008,
		"gateway":   8000,
		"auth":      8009,
		"admin":     8010,
	}

	if port, ok := ports[serviceName]; ok {
		return port
	}
	return 8080
}

// Assign unique default gRPC port for each service
func getDefaultGRPCPort(serviceName string) int {
	ports := map[string]int{
		"user":      9001,
		"product":   9002,
		"inventory": 9003,
		"order":     9004,
		"payment":   9005,
		"marketing": 9006,
		"cms":       9007,
		"shipping":  9008,
		"gateway":   9000,
		"auth":      9009,
		"admin":     9010,
	}

	if port, ok := ports[serviceName]; ok {
		return port
	}
	return 9090
}
