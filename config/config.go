package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	FirebaseAPIKey string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file")
	}

	return &Config{
		FirebaseAPIKey: os.Getenv("FIREBASE_API_KEY"),
	}
}
