package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BrapiToken string
}

// Load reads configuration from the .env file and environment variables.
// Environment variables take precedence over the .env file.
func Load() (*Config, error) {
	// Ignore error — .env file is optional (env vars may be set directly)
	godotenv.Load(".env")

	token := os.Getenv("BRAPI_TOKEN")
	if token == "" {
		return nil, errors.New("BRAPI_TOKEN is required: set it in .env or as an environment variable")
	}

	return &Config{BrapiToken: token}, nil
}
