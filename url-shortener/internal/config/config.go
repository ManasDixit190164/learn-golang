package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	Port               string
	BaseURL            string
	DatabaseURL        string
	JWTAccessSecret    string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	ShortCodeLength    int
}

func Load() (Config, error) {
	_ = godotenv.Load()

	accessExpiry, err := time.ParseDuration(getEnv("ACCESS_TOKEN_EXPIRY", "15m"))
	if err != nil {
		return Config{}, err
	}

	refreshExpiry, err := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))
	if err != nil {
		return Config{}, err
	}

	shortCodeLength, err := strconv.Atoi(getEnv("SHORT_CODE_LENGTH", "7"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		BaseURL:            getEnv("BASE_URL", "http://localhost:8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTAccessSecret:    os.Getenv("JWT_ACCESS_SECRET"),
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
		ShortCodeLength:    shortCodeLength,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	if cfg.JWTAccessSecret == "" || cfg.JWTAccessSecret == "change-this-access-secret" {
		return Config{}, errors.New("JWT_ACCESS_SECRET must be set to a strong secret")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
