package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr         string
	POSTGRES_URI string
	FRONT_URI    string

	JWTSecretAccess  string
	JWTSecretRefresh string
	JWTAccessTTL     time.Duration
	JWTRefreshTTL    time.Duration
	RedisAddr        string
	RedisPassword    string

	SMTP_PORT string
	SMTP_HOST string
	SMTP_USER string
	SMTP_PASS string

	SWAGGER_LOGIN string
	SWAGGER_PASS  string

	PATH_IMAGE string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Addr:         getEnv("PUBLIC_HOST", "http://localhost"),
		POSTGRES_URI: getEnv("POSTGRES_URI", "postgres://test:test@localhost:5432/test"),
		FRONT_URI:    getEnv("FRONT_URI", "http://localhost:3000"),

		JWTSecretAccess:  getEnv("JWTSecretAccess", "secretAccess"),
		JWTSecretRefresh: getEnv("JWTSecretRefresh", "secretRefresh"),
		JWTAccessTTL:     15 * time.Minute,
		JWTRefreshTTL:    30 * 24 * time.Hour,
		RedisAddr:        getEnv("RedisAddr", "localhost:6379"),
		RedisPassword:    getEnv("RedisPassword", ""),

		SMTP_PORT: getEnv("SMTP_PORT", "587"),
		SMTP_HOST: getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTP_USER: getEnv("SMTP_USER", "firulvv@gmail.com"),
		SMTP_PASS: getEnv("SMTP_PASS", "rkvt hdki wksz phnp"),

		SWAGGER_LOGIN: getEnv("SMTP_PASS", "admin"),
		SWAGGER_PASS:  getEnv("SMTP_PASS", "12345"),

		PATH_IMAGE: getEnv("PATH_IMAGE", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
