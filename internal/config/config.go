package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	APIBaseURL string
	APITimeout int
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}

	timeout, err := strconv.Atoi(os.Getenv("API_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid API_TIMEOUT: %v", err)
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     port,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		APIBaseURL: os.Getenv("API_BASE_URL"),
		APITimeout: timeout,
	}, nil
}
