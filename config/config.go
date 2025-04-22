package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr         string
	POSTGRES_URI string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Addr:         getEnv("PUBLIC_HOST", "http://localhost"),
		POSTGRES_URI: getEnv("POSTGRES_URI", "postgres://test:test@localhost:5432/test"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
