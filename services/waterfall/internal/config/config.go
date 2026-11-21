package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
	PoolSize int
}

// RedisConfig contains Redis connection settings
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// AuthConfig contains authentication settings
type AuthConfig struct {
	Enabled       bool
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshSecret string
	RefreshExpiry time.Duration
	Issuer        string
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Port            string
	Mode            string
	ShutdownTimeout time.Duration
}

// WaterfallConfig contains Waterfall-specific settings
type WaterfallConfig struct {
	WorkerCount           int
	QueueWorkerCount      int
	SchedulerInterval     time.Duration
	CapacityCheckInterval time.Duration
	FailureCheckInterval  time.Duration
	CostOptimizerInterval time.Duration
	MaxRetries            int
	RetryDelay            time.Duration
	EnableAI              bool
	AICheckInterval       time.Duration
}

// Config contains all application configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Auth      AuthConfig
	Waterfall WaterfallConfig
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	// Server configuration
	cfg.Server.Port = getEnv("SERVER_PORT", "8081")
	cfg.Server.Mode = getEnv("SERVER_MODE", "release")
	cfg.Server.ShutdownTimeout = getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second)

	// Database configuration
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getIntEnv("DB_PORT", 5432)
	cfg.Database.Username = getEnv("DB_USERNAME", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.Database = getEnv("DB_NAME", "brlia_waterfall")
	cfg.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")
	cfg.Database.PoolSize = getIntEnv("DB_POOL_SIZE", 10)

	// Redis configuration
	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getIntEnv("REDIS_PORT", 6379)
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getIntEnv("REDIS_DB", 1)

	// Auth configuration
	cfg.Auth.Enabled = getBoolEnv("AUTH_ENABLED", true)
	cfg.Auth.JWTSecret = getEnv("JWT_SECRET", "brlia-waterfall-secret")
	cfg.Auth.JWTExpiry = getDurationEnv("JWT_EXPIRY", 15*time.Minute)
	cfg.Auth.RefreshSecret = getEnv("REFRESH_SECRET", "brlia-refresh-secret")
	cfg.Auth.RefreshExpiry = getDurationEnv("REFRESH_EXPIRY", 7*24*time.Hour)
	cfg.Auth.Issuer = getEnv("JWT_ISSUER", "brlia-waterfall-api")

	// Waterfall configuration
	cfg.Waterfall.WorkerCount = getIntEnv("WORKER_COUNT", 10)
	cfg.Waterfall.QueueWorkerCount = getIntEnv("QUEUE_WORKER_COUNT", 5)
	cfg.Waterfall.SchedulerInterval = getDurationEnv("SCHEDULER_INTERVAL", 1*time.Second)
	cfg.Waterfall.CapacityCheckInterval = getDurationEnv("CAPACITY_CHECK_INTERVAL", 5*time.Second)
	cfg.Waterfall.FailureCheckInterval = getDurationEnv("FAILURE_CHECK_INTERVAL", 10*time.Second)
	cfg.Waterfall.CostOptimizerInterval = getDurationEnv("COST_OPTIMIZER_INTERVAL", 30*time.Second)
	cfg.Waterfall.MaxRetries = getIntEnv("MAX_RETRIES", 3)
	cfg.Waterfall.RetryDelay = getDurationEnv("RETRY_DELAY", 1*time.Second)
	cfg.Waterfall.EnableAI = getBoolEnv("ENABLE_AI", false)
	cfg.Waterfall.AICheckInterval = getDurationEnv("AI_CHECK_INTERVAL", 1*time.Minute)

	// Validate required configurations
	if cfg.Auth.Enabled && cfg.Auth.JWTSecret == "brlia-waterfall-secret" {
		return nil, fmt.Errorf("JWT_SECRET must be configured in production")
	}

	return cfg, nil
}

// getEnv retrieves an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv retrieves an integer environment variable with a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getBoolEnv retrieves a boolean environment variable with a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// getDurationEnv retrieves a duration environment variable with a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
