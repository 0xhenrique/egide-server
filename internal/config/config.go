package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	ServerPort  int
	DatabaseURL string
	GitHubOAuth struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
		Scopes       []string
	}
	JWTSecret string
}

func New() (*Config, error) {
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT")
	}

	dbURL := getEnv("DATABASE_URL", "./data.db")
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	githubClientID := getEnv("GITHUB_CLIENT_ID", "")
	githubClientSecret := getEnv("GITHUB_CLIENT_SECRET", "")
	if githubClientID == "" || githubClientSecret == "" {
		return nil, errors.New("GitHub OAuth credentials are required")
	}

	cfg := &Config{
		ServerPort:  port,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
	}

	cfg.GitHubOAuth.ClientID = githubClientID
	cfg.GitHubOAuth.ClientSecret = githubClientSecret
	cfg.GitHubOAuth.RedirectURL = getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/callback")
	cfg.GitHubOAuth.Scopes = []string{"user:email"}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
