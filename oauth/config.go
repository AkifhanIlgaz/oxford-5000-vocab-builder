package oauth

import "fmt"

type OAuthHandlers struct {
	Google *GoogleOAuth
}

func NewOAuthHandlers() (*OAuthHandlers, error) {
	google, err := NewGoogleOauth()
	if err != nil {
		return nil, fmt.Errorf("new google oauth config: %w", err)
	}

	return &OAuthHandlers{
		Google: google,
	}, nil
}
