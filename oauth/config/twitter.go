package config

import (
	"errors"
	"os"
)

type TwitterOAuth struct {
	ClientKey    string
	ClientSecret string
	AuthUrl      string
	TokenUrl     string
	UserUrl      string
}

func NewTwitterOauth() (*TwitterOAuth, error) {
	key := os.Getenv("TWITTER_CLIENT_KEY")
	if key == "" {
		return nil, errors.New("twitter client key required")
	}

	secret := os.Getenv("TWITTER_CLIENT_KEY")
	if secret == "" {
		return nil, errors.New("twitter secret key required")
	}

	return &TwitterOAuth{
		ClientKey:    key,
		ClientSecret: secret,
		AuthUrl:      "https://api.twitter.com/oauth/authorize",
		TokenUrl:     "https://api.twitter.com/oauth2/token",
		UserUrl:      "https://api.twitter.com/2/users/me",
	}, nil
}
