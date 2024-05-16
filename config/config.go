package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	ServerPort string

	DSN string

	LogLevel string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	return &Config{
		ServerHost: os.Getenv("SERVER_HOST"),
		ServerPort: os.Getenv("SERVER_PORT"),

		DSN: os.Getenv("DSN"),

		LogLevel: os.Getenv("LOG_LEVEL"),
	}, nil
}
