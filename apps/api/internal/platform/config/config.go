package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	APIAddr         string
	DatabaseURL     string
	SessionDuration time.Duration
	SessionSecret   string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}

	return Config{
		AppEnv:          stringOrDefault("APP_ENV", "development"),
		APIAddr:         stringOrDefault("API_ADDR", ":8080"),
		DatabaseURL:     stringOrDefault("DATABASE_URL", "postgres://user:pass@localhost:5432/gated_db?sslmode=disable"),
		SessionDuration: durationOrDefault("SESSION_DURATION", 7*24*time.Hour),
		SessionSecret:   stringOrDefault("SESSION_SECRET", "dev-session-secret-change-me"),
	}, nil
}

func (c Config) IsDevelopment() bool {
	return strings.EqualFold(c.AppEnv, "development")
}

func stringOrDefault(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}

	return value
}

func durationOrDefault(key string, fallback time.Duration) time.Duration {
	value := stringOrDefault(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}
