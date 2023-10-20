package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/AkifhanIlgaz/vocab-builder/rand"
)

type GoogleOAuth struct {
	ClientKey    string
	ClientSecret string
	AuthUrl      string
	TokenUrl     string
	UserUrl      string
}

func NewGoogleOauth() (*GoogleOAuth, error) {
	key := os.Getenv("GOOGLE_CLIENT_ID")
	if key == "" {
		return nil, errors.New("google client key required")
	}

	secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if secret == "" {
		return nil, errors.New("google secret key required")
	}

	// TODO: Add redirect uri
	return &GoogleOAuth{
		ClientKey:    key,
		ClientSecret: secret,
		AuthUrl:      "https://accounts.google.com/o/oauth2/auth",
		TokenUrl:     "https://oauth2.googleapis.com/token",
		UserUrl:      "https://www.googleapis.com/oauth2/v3/userinfo",
	}, nil
}

func (google *GoogleOAuth) Signin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("client_id", google.ClientKey)
	query.Set("access_type", "offline")
	query.Set("response_type", "code")
	query.Set("scope", "email")
	query.Set("redirect_uri", "http://localhost:3000/auth/google/callback")
	if state, err := rand.String(32); err == nil {
		query.Set("state", state)
	}

	url := fmt.Sprintf("%s?%s", google.AuthUrl, query.Encode())
	fmt.Println(url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (google *GoogleOAuth) Callback(w http.ResponseWriter, r *http.Request) {

}
