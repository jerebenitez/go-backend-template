package utils

import (
	"os"

	// this will automatically load your .env file:
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	PublicHost string
	Port string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		PublicHost: getEnv("HOST", "localhost"),
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
