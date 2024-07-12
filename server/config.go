package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ApiKey      string
	AdminKey    string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	url := os.Getenv("SERVICE_URL")
	if url == "" {
		panic("SERVICE_URL is not set")
	}

	apiKey := os.Getenv("SERVICE_KEY")
	if apiKey == "" {
		panic("SERVICE_KEY is not set")
	}

	adminKey := os.Getenv("ADMIN_KEY")
	if adminKey == "" {
		panic("ADMIN_KEY is not set")
	}

	return Config{
		DatabaseURL: url,
		ApiKey:      apiKey,
		AdminKey:    adminKey,
	}
}
