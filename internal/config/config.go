package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	return godotenv.Load()
}

func Get(key string) string {
	return os.Getenv(key)
}
