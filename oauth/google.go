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

const (
	invalidRequest     = "invalid_request"
	invalidCredentials = "Invalid Credentials"
)

type tokenInfo struct {
	Provider     Provider `json:"provider"`
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	Scope        string   `json:"scope"`
	IdToken      string   `json:"id_token"`
}

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
			RedirectURL:  "http://localhost:3000/auth/signin/google/callback",
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"}},
	}, nil

}

func (google *GoogleOAuth) Signin(w http.ResponseWriter, r *http.Request) {
	authUrl := google.AuthCodeURL("random state", oauth2.AccessTypeOffline)

	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
}

func (google *GoogleOAuth) Callback(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")

	token, err := google.Exchange(context.TODO(), code, oauth2.AccessTypeOffline)
	if err != nil {
		fmt.Println("err", err)
	}

	json.NewEncoder(w).Encode(token)

}

func (google *GoogleOAuth) AccessTokenMiddleware(w http.ResponseWriter, r *http.Request) {
	var token oauth2.Token

	b, _ := io.ReadAll(r.Body)

	json.Unmarshal(b, &token)

	res, _ := google.Client(context.TODO(), &token).Get(google.UserUrl)

	body, _ := io.ReadAll(res.Body)
	var respBody map[string]string
	json.Unmarshal(body, &respBody)

	// TODO: Set r.Context with uid

	json.NewEncoder(w).Encode(&respBody)

}

func (google *GoogleOAuth) GenerateAccessTokenWithRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refresh_token")

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
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, _ := io.ReadAll(res.Body)

	var data tokenInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&data)

}

func (google *GoogleOAuth) createTokenRequest(r *http.Request) (*http.Request, error) {
	code := r.URL.Query().Get("code")

	requestBodyMap := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     google.ClientID,
		"client_secret": google.ClientSecret,
		"code":          code,
		"redirect_uri":  google.RedirectURL,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest(
		"POST",
		google.Endpoint.TokenURL,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}
