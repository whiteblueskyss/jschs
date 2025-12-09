package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration.
type Config struct {
	DatabaseURL string
	ServerAddr  string
}

// Load reads configuration from environment variables.
func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it:", err)
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ServerAddr:  os.Getenv("SERVER_ADDR"),
	}

	return cfg
}
