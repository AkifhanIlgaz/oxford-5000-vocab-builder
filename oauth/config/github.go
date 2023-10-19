package config

import (
	"errors"
	"os"
)

type GithubOAuthConfig struct {
	ClientKey    string
	ClientSecret string
	AuthUrl      string
	TokenUrl     string
	UserUrl      string
}

func NewGithubConfig() (*GithubOAuthConfig, error) {
	key := os.Getenv("GITHUB_CLIENT_KEY")
	if key == "" {
		return nil, errors.New("github client key required")
	}

	secret := os.Getenv("GITHUB_CLIENT_KEY")
	if secret == "" {
		return nil, errors.New("github secret key required")
	}

	return &GithubOAuthConfig{
		ClientKey:    key,
		ClientSecret: secret,
		AuthUrl:      "https://github.com/login/oauth/authorize",
		TokenUrl:     "https://github.com/login/oauth/access_token",
		UserUrl:      "https://api.github.com/user",
	}, nil
}
