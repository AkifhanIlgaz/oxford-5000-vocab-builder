package oauth

import "fmt"

type OAuthHandlers struct {
	Github  *GithubOAuth
	Google  *GoogleOAuth
	Twitter *TwitterOAuth
}

func NewOAuthHandlers() (*OAuthHandlers, error) {
	github, err := NewGithubOAuth()
	if err != nil {
		return nil, fmt.Errorf("new github oauth config: %w", err)
	}

	google, err := NewGoogleOauth()
	if err != nil {
		return nil, fmt.Errorf("new google oauth config: %w", err)
	}

	// TODO: Add twitter

	// twitter, err := NewTwitterOauth()
	// if err != nil {
	// 	return nil, fmt.Errorf("new twitter oauth config: %w", err)
	// }

	// TODO: Add instagram

	return &OAuthHandlers{
		Google: google,
		// Twitter: twitter,
		Github: github,
	}, nil
}
