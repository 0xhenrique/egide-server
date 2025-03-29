package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	JWTSecret       string
	TokenExpiration time.Duration
	DBPath          string
}

func Load() *Config {
	config := &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
		DBPath:    getEnv("DB_PATH", "./egide.db"),
	}

	expirationHours, err := strconv.Atoi(getEnv("TOKEN_EXPIRATION_HOURS", "24"))
	if err != nil {
		expirationHours = 24
	}
	config.TokenExpiration = time.Duration(expirationHours) * time.Hour

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
