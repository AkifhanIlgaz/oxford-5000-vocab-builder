package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	query.Set("redirect_uri", "http://localhost:3000/auth/signin/google/callback")
	// ? How to use state
	if state, err := rand.String(32); err == nil {
		query.Set("state", state)
	}

	url := fmt.Sprintf("%s?%s", google.AuthUrl, query.Encode())

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (google *GoogleOAuth) Callback(w http.ResponseWriter, r *http.Request) {
	// TODO: Check state
	code := r.URL.Query().Get("code")

	requestBodyMap := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     google.ClientKey,
		"client_secret": google.ClientSecret,
		"code":          code,
		"redirect_uri":  "http://localhost:3000/auth/signin/google/callback",
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest(
		"POST",
		google.TokenUrl,
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respbody))
	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var respBody map[any]any
	json.Unmarshal(respbody, &respBody)

	fmt.Fprint(w, respBody)

}
