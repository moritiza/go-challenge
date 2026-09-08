package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration.
type Config struct {
	Server ServerConfig
	Redis  RedisConfig
	Domain DomainConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// RedisConfig holds connection settings for Redis.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// DomainConfig holds business-rule settings.
type DomainConfig struct {
	MembershipTTL time.Duration
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

	shutdownTimeout, err := time.ParseDuration(getEnv("SERVER_SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_SHUTDOWN_TIMEOUT: %w", err)
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	membershipTTL, err := time.ParseDuration(getEnv("MEMBERSHIP_TTL", "336h")) // 14 days
	if err != nil {
		return nil, fmt.Errorf("invalid MEMBERSHIP_TTL: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			ShutdownTimeout: shutdownTimeout,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Domain: DomainConfig{
			MembershipTTL: membershipTTL,
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
