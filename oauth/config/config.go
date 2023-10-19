package config

import "fmt"

type OAuthConfig struct {
	Github  *GithubOAuth
	Google  *GoogleOAuth
	Twitter *TwitterOAuth
}

func NewOAuthConfig() (*OAuthConfig, error) {
	github, err := NewGithubOAuth()
	if err != nil {
		return nil, fmt.Errorf("new oauth config: %w", err)
	}

	google, err := NewGoogleOauth()
	if err != nil {
		return nil, fmt.Errorf("new oauth config: %w", err)
	}

	twitter, err := NewTwitterOauth()
	if err != nil {
		return nil, fmt.Errorf("new oauth config: %w", err)
	}

	return &OAuthConfig{
		Google:  google,
		Twitter: twitter,
		Github:  github,
	}, nil
}
