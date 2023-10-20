package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	RedirectUrl  string
}

type tokenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	IdToken      string `json:"id_token"`
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

	return &GoogleOAuth{
		ClientKey:    key,
		ClientSecret: secret,
		AuthUrl:      "https://accounts.google.com/o/oauth2/auth",
		TokenUrl:     "https://oauth2.googleapis.com/token",
		UserUrl:      "https://www.googleapis.com/oauth2/v3/userinfo",
		RedirectUrl:  "http://localhost:3000/auth/signin/google/callback",
	}, nil
}

func (google *GoogleOAuth) Signin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("client_id", google.ClientKey)
	query.Set("access_type", "offline")
	query.Set("response_type", "code")
	query.Set("scope", "email")
	query.Set("redirect_uri", google.RedirectUrl)
	// ? How to use state
	if state, err := rand.String(32); err == nil {
		query.Set("state", state)
	}

	url := fmt.Sprintf("%s?%s", google.AuthUrl, query.Encode())

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (google *GoogleOAuth) Callback(w http.ResponseWriter, r *http.Request) {
	req, err := google.createTokenRequest(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	tokenInfo, err := google.exchangeCodeForTokens(req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(tokenInfo)
}

func (google *GoogleOAuth) VerifyAccessToken(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	fmt.Println(accessToken)

	req, _ := http.NewRequest("GET", google.UserUrl, nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

func (google *GoogleOAuth) createTokenRequest(r *http.Request) (*http.Request, error) {
	code := r.URL.Query().Get("code")

	requestBodyMap := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     google.ClientKey,
		"client_secret": google.ClientSecret,
		"code":          code,
		"redirect_uri":  google.RedirectUrl,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest(
		"POST",
		google.TokenUrl,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (google *GoogleOAuth) exchangeCodeForTokens(r *http.Request) (*tokenInfo, error) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("exchange code for tokens: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("exchange code for tokens: %w", err)
	}

	var data tokenInfo
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, fmt.Errorf("exchange code for tokens: %w", err)
	}

	return &data, nil
}
