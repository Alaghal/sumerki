package config

import (
	"fmt"
	"os"
)

const defaultBackendPort = "8080"

type Config struct {
	DatabaseURL string
	JWTSecret   string
	BackendPort string
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		BackendPort: os.Getenv("BACKEND_PORT"),
	}

	if cfg.BackendPort == "" {
		cfg.BackendPort = defaultBackendPort
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}
