package oauth

import (
	"fmt"
)

func Setup() (*OAuthHandlers, error) {
	oauth, err := NewOAuthHandlers()
	if err != nil {
		return nil, fmt.Errorf("oauth setup: %w", err)
	}

	return oauth, err
}
