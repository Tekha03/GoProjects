package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv string
	Port   string

	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	JWT struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	cfg.AppEnv = getEnv("APP_ENV", "development")
	cfg.Port = getEnv("PORT", "8080")

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.DB.Name = getEnv("DB_NAME", "marketplace")
	cfg.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	accessMinutes, err := strconv.Atoi(getEnv("JWT_ACCESS_TTL_MIN", "15"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL_MIN")
	}
	cfg.JWT.AccessTTL = time.Duration(accessMinutes) * time.Minute

	refreshHours, err := strconv.Atoi(getEnv("JWT_REFRESH_TTL_HOURS", "168"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL_HOURS")
	}
	cfg.JWT.RefreshTTL = time.Duration(refreshHours) * time.Hour

	return cfg, nil
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
