package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPAddr    string
	SwaggerHost string
	DatabaseURL string
	JWTSecret   string
	JWTTTL      time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		SwaggerHost: getEnv("SWAGGER_HOST", "localhost:8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://filetemp:filetemp@localhost:5432/filetemp?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-jwt-secret"),
	}

	ttlRaw := getEnv("JWT_TTL", "24h")
	ttl, err := time.ParseDuration(ttlRaw)
	if err != nil {
		return Config{}, fmt.Errorf("JWT_TTL: %w", err)
	}
	if ttl <= 0 {
		return Config{}, fmt.Errorf("JWT_TTL must be positive")
	}
	cfg.JWTTTL = ttl

	if cfg.JWTSecret == "" || cfg.JWTSecret == "change-me-jwt-secret" {
		// allow default in local/dev; warn is caller's concern
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
