package oauth

import (
	"errors"
	"os"
)

type GithubOAuth struct {
	ClientKey    string
	ClientSecret string
	AuthUrl      string
	TokenUrl     string
	UserUrl      string
}

func NewGithubOAuth() (*GithubOAuth, error) {
	key := os.Getenv("GITHUB_CLIENT_ID")
	if key == "" {
		return nil, errors.New("github client key required")
	}

	secret := os.Getenv("GITHUB_CLIENT_SECRET")
	if secret == "" {
		return nil, errors.New("github secret key required")
	}

	return &GithubOAuth{
		ClientKey:    key,
		ClientSecret: secret,
		AuthUrl:      "https://github.com/login/oauth/authorize",
		TokenUrl:     "https://github.com/login/oauth/access_token",
		UserUrl:      "https://api.github.com/user",
	}, nil
}
