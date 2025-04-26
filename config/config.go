package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr         string
	POSTGRES_URI string

	JWTSecretAccess  string
	JWTSecretRefresh string
	JWTAccessTTL     time.Duration
	JWTRefreshTTL    time.Duration
	RedisAddr        string
	RedisPassword    string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Addr:         getEnv("PUBLIC_HOST", "http://localhost"),
		POSTGRES_URI: getEnv("POSTGRES_URI", "postgres://test:test@localhost:5432/test"),

		JWTSecretAccess:  "secretAccess",
		JWTSecretRefresh: "secretRefresh",
		JWTAccessTTL:     15 * time.Minute,
		JWTRefreshTTL:    30 * 24 * time.Hour,
		RedisAddr:        "localhost:6379",
		RedisPassword:    "",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
