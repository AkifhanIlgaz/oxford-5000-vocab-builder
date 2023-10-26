package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	UserUrl string
	oauth2.Config
}

func NewGoogleOauth() (*GoogleOAuth, error) {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	if clientId == "" {
		return nil, errors.New("google client key required")
	}

	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, errors.New("google secret key required")
	}

	return &GoogleOAuth{
		UserUrl: "https://www.googleapis.com/oauth2/v3/tokeninfo",
		Config: oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  "http://localhost:3000/auth/google/callback",
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"}},
	}, nil

}

func (google *GoogleOAuth) Signin(w http.ResponseWriter, r *http.Request) {
	authUrl := google.AuthCodeURL("", oauth2.AccessTypeOffline)

	// TODO: Generate random string for state
	// TODO: Store it on cache ?

	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
}

func (google *GoogleOAuth) Callback(w http.ResponseWriter, r *http.Request) {
	// TODO: Check state

	token, err := google.Exchange(context.TODO(), r.URL.Query().Get("code"))
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

func (google *GoogleOAuth) GenerateAccessTokenWithRefreshToken(refreshToken string) (*oauth2.Token, error) {
	requestBodyMap := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     google.ClientID,
		"client_secret": google.ClientSecret,
		"refresh_token": refreshToken,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest(
		"POST",
		google.Endpoint.TokenURL,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("generate access token with refresh: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("generate access token with refresh: %w", err)
	}

	body, _ := io.ReadAll(res.Body)

	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil

}
