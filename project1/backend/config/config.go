package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Database: DatabaseConfig{
			DSN: getEnv("DATABASE_DSN", "./taskmanager.db"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			ExpirationHours: 24,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
