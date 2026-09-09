package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration.
type Config struct {
	Server      ServerConfig
	Redis       RedisConfig
	Domain      DomainConfig
	Maintenance MaintenanceConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// RedisConfig holds connection settings for Redis.
type RedisConfig struct {
	Addr        string
	Password    string
	DB          int
	PingTimeout time.Duration
}

// DomainConfig holds business-rule settings.
type DomainConfig struct {
	MembershipTTL time.Duration
}

type MaintenanceConfig struct {
	CleanupInterval  time.Duration
	CleanupBatchSize int64
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	readTimeout, err := time.ParseDuration(getEnv("SERVER_READ_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := time.ParseDuration(getEnv("SERVER_WRITE_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
	}

	idleTimeout, err := time.ParseDuration(getEnv("SERVER_IDLE_TIMEOUT", "60s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_IDLE_TIMEOUT: %w", err)
	}

	shutdownTimeout, err := time.ParseDuration(getEnv("SERVER_SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_SHUTDOWN_TIMEOUT: %w", err)
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	pingTimeout, err := time.ParseDuration(getEnv("REDIS_PING_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PING_TIMEOUT: %w", err)
	}

	membershipTTL, err := time.ParseDuration(getEnv("MEMBERSHIP_TTL", "336h")) // 14 days
	if err != nil {
		return nil, fmt.Errorf("invalid MEMBERSHIP_TTL: %w", err)
	}

	cleanupInterval, err := time.ParseDuration(getEnv("CLEANUP_INTERVAL", "10m"))
	if err != nil {
		return nil, fmt.Errorf("invalid CLEANUP_INTERVAL: %w", err)
	}

	cleanupBatchSize, err := strconv.ParseInt(getEnv("CLEANUP_BATCH_SIZE", "1000"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid CLEANUP_BATCH_SIZE: %w", err)
	}

	if cleanupBatchSize <= 0 {
		return nil, fmt.Errorf("invalid CLEANUP_BATCH_SIZE: must be greater than 0")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			IdleTimeout:     idleTimeout,
			ShutdownTimeout: shutdownTimeout,
		},
		Redis: RedisConfig{
			Addr:        getEnv("REDIS_ADDR", "localhost:6379"),
			Password:    getEnv("REDIS_PASSWORD", ""),
			DB:          redisDB,
			PingTimeout: pingTimeout,
		},
		Domain: DomainConfig{
			MembershipTTL: membershipTTL,
		},
		Maintenance: MaintenanceConfig{
			CleanupInterval:  cleanupInterval,
			CleanupBatchSize: cleanupBatchSize,
		},
	}

	return cfg, nil
}

// getEnv retrieves an environment variable with a fallback value.
func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}

	return fallback
}
