package utils

import (
	"fmt"
	"os"

	// this will automatically load your .env file:
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	PublicHost string
	Port       string
	DB         DbConfig
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		PublicHost: getEnv("HOST", ""),
		Port:       getEnv("PORT", "5000"),
		DB: DbConfig{
			DSN:      getEnv("DB_DSN", ""),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "admin"),
			Name:     getEnv("DB_NAME", "postgres"),
			Path: fmt.Sprintf(
				"%s:%s",
				getEnv("DB_HOST", "localhost"),
				getEnv("DB_PORT", "5432"),
			),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
