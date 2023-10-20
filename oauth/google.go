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

const (
	invalidRequest     = "invalid_request"
	invalidCredentials = "Invalid Credentials"
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
		TokenUrl:     "https://www.googleapis.com/oauth2/v4/token",
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

// middleware
func (google *GoogleOAuth) VerifyAccessToken(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")

	req, _ := http.NewRequest("GET", google.UserUrl, nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("err", err)
	}

	body, _ := io.ReadAll(res.Body)
	var respBody map[string]string
	json.Unmarshal(body, &respBody)

	if isInvalidRequest(respBody) {
		// TODO: http return invalid token error
		// Client will go to /refresh endpoint to obtain new access token
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&respBody)
		return
	}

	json.NewEncoder(w).Encode(&respBody)

}

func isInvalidRequest(respBody map[string]string) bool {
	return respBody["error"] == invalidRequest && respBody["error_description"] == invalidCredentials
}

func (google *GoogleOAuth) GenerateAccessTokenWithRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refresh_token")

	requestBodyMap := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     google.ClientKey,
		"client_secret": google.ClientSecret,
		"refresh_token": refreshToken,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest(
		"POST",
		google.TokenUrl,
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
