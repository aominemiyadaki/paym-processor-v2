package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string
	Port  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DBURL: dbURL,
		Port:  port,
	}, nil
}
