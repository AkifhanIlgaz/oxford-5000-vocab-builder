package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type GithubOAuth struct {
	oauth2.Config
}

func NewGithubOAuth() (*GithubOAuth, error) {
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	if clientId == "" {
		return nil, errors.New("github client id required")
	}

	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, errors.New("github client required")
	}

	return &GithubOAuth{
		Config: oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  "http://localhost:3000/auth/github/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		},
	}, nil

}

func (github *GithubOAuth) Signin(w http.ResponseWriter, r *http.Request) {
	authUrl := github.AuthCodeURL("", oauth2.AccessTypeOffline)

	// TODO: Generate random string for state
	// TODO: Store it on cache ?

	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
}

func (github *GithubOAuth) Callback(w http.ResponseWriter, r *http.Request) {
	token, err := github.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "http://localhost:8100/test/error", http.StatusUnauthorized)
		return
	}

	info, err := json.Marshal(&token)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, "http://localhost:8100/test?info="+string(info), http.StatusTemporaryRedirect)

}
