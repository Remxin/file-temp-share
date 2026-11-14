package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr         string
	SwaggerHost      string
	DatabaseURL      string
	JWTSecret        string
	JWTTTL           time.Duration
	RabbitMQURL      string
	AppPublicURL     string
	GmailAppLogin    string
	GmailAppPassword string

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool
	UploadMaxBytes int64
}

// Load reads optional .env / .env.local files then environment variables.
func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		HTTPAddr:         getEnv("HTTP_ADDR", ":8080"),
		SwaggerHost:      getEnv("SWAGGER_HOST", "localhost:8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://filetemp:filetemp@localhost:5432/filetemp?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "change-me-jwt-secret"),
		RabbitMQURL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		AppPublicURL:     stringsTrimRightSlash(getEnv("APP_PUBLIC_URL", "http://localhost:3000")),
		GmailAppLogin:    os.Getenv("GMAIL_APP_LOGIN"),
		GmailAppPassword: os.Getenv("GMAIL_APP_PASSWORD"),
		MinIOEndpoint:    getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey:   getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey:   getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:      getEnv("MINIO_BUCKET", "transfers"),
		MinIOUseSSL:      parseBool(getEnv("MINIO_USE_SSL", "false")),
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

	maxBytes, err := strconv.ParseInt(getEnv("UPLOAD_MAX_BYTES", "536870912"), 10, 64)
	if err != nil || maxBytes <= 0 {
		return Config{}, fmt.Errorf("UPLOAD_MAX_BYTES must be a positive integer")
	}
	cfg.UploadMaxBytes = maxBytes

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.RabbitMQURL == "" {
		return Config{}, fmt.Errorf("RABBITMQ_URL is required")
	}
	if cfg.AppPublicURL == "" {
		return Config{}, fmt.Errorf("APP_PUBLIC_URL is required")
	}
	if cfg.MinIOEndpoint == "" {
		return Config{}, fmt.Errorf("MINIO_ENDPOINT is required")
	}
	if cfg.MinIOAccessKey == "" || cfg.MinIOSecretKey == "" {
		return Config{}, fmt.Errorf("MINIO_ACCESS_KEY and MINIO_SECRET_KEY are required")
	}
	if cfg.MinIOBucket == "" {
		return Config{}, fmt.Errorf("MINIO_BUCKET is required")
	}
	return cfg, nil
}

// RequireGmail ensures SMTP credentials are present (worker).
func (c Config) RequireGmail() error {
	if c.GmailAppLogin == "" {
		return fmt.Errorf("GMAIL_APP_LOGIN is required")
	}
	if c.GmailAppPassword == "" {
		return fmt.Errorf("GMAIL_APP_PASSWORD is required")
	}
	return nil
}

func loadDotEnv() {
	// Prefer repo-root files when running from backend/ via Taskfile.
	candidates := []string{
		".env",
		".env.local",
		filepath.Join("..", ".env"),
		filepath.Join("..", ".env.local"),
	}
	var base []string
	var local []string
	for _, p := range candidates {
		if fileExists(p) {
			if filepath.Base(p) == ".env.local" {
				local = append(local, p)
			} else {
				base = append(base, p)
			}
		}
	}
	if len(base) > 0 {
		_ = godotenv.Load(base...)
	}
	if len(local) > 0 {
		_ = godotenv.Overload(local...)
	}
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func stringsTrimRightSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
