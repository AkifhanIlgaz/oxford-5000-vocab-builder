package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GithubOAuth struct {
	ClientKey    string
	ClientSecret string
	AuthUrl      string
	TokenUrl     string
	UserUrl      string
	RedirectUrl  string
}

func NewGithubOAuth() (*GithubOAuth, error) {
	key := os.Getenv("GITHUB_CLIENT_ID")
	if key == "" {
		return nil, errors.New("github client key required")
	}

	secret := os.Getenv("GITHUB_CLIENT_SECRET")
	if secret == "" {
		return nil, errors.New("github secret key required")
	}

	return &GithubOAuth{
		ClientKey:    key,
		ClientSecret: secret,
		AuthUrl:      "https://github.com/login/oauth/authorize",
		TokenUrl:     "https://github.com/login/oauth/access_token",
		UserUrl:      "https://api.github.com/user",
		RedirectUrl:  "http://localhost:3000/auth/signin/github/callback",
	}, nil
}

func (github *GithubOAuth) Signin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	query.Set("client_id", github.ClientKey)
	query.Set("redirect_uri", github.RedirectUrl)

	url := fmt.Sprintf("%s?%s", github.AuthUrl, query.Encode())

	fmt.Println(url)
	w.Header().Set("Origin", r.Header.Get("Origin")+"/")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (github *GithubOAuth) Callback(w http.ResponseWriter, r *http.Request) {
	req, err := github.createTokenRequest(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	github.exchangeCodeForTokens(req)

	// tokenInfo, err := github.exchangeCodeForTokens(req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, "Something went wrong", http.StatusInternalServerError)
	// 	return
	// }
	// tokenInfo.Provider = ProviderGoogle

	// json.NewEncoder(w).Encode(tokenInfo)
}

func (github *GithubOAuth) createTokenRequest(r *http.Request) (*http.Request, error) {
	code := r.URL.Query().Get("code")

	requestBodyMap := map[string]string{
		"client_id":     github.ClientKey,
		"client_secret": github.ClientSecret,
		"code":          code,
		"redirect_uri":  github.RedirectUrl,
	}

	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest(
		"POST",
		github.TokenUrl,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (github *GithubOAuth) exchangeCodeForTokens(req *http.Request) {
	resp, _ := http.DefaultClient.Do(req)

	respBody, _ := io.ReadAll(resp.Body)

	fmt.Println(string(respBody))

	var data map[string]any
	json.Unmarshal(respBody, &data)

	fmt.Println(data)

	// return &data, nil
}
