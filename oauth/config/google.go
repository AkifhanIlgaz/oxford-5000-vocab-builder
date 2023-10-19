package config

import (
	"errors"
	"os"
)

type GoogleOAuth struct {
	ClientKey    string
	ClientSecret string
	AuthUrl      string
	TokenUrl     string
	UserUrl      string
}

func NewGoogleOauth() (*GoogleOAuth, error) {
	key := os.Getenv("GOOGLE_CLIENT_KEY")
	if key == "" {
		return nil, errors.New("google client key required")
	}

	secret := os.Getenv("GOOGLE_CLIENT_KEY")
	if secret == "" {
		return nil, errors.New("google secret key required")
	}

	return &GoogleOAuth{
		ClientKey:    key,
		ClientSecret: secret,
		AuthUrl:      "https://accounts.google.com/o/oauth2/auth",
		TokenUrl:     "https://oauth2.googleapis.com/token",
		UserUrl:      "https://www.googleapis.com/oauth2/v3/userinfo",
	}, nil
}
