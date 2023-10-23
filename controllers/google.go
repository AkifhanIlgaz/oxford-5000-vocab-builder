package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AkifhanIlgaz/vocab-builder/errors"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleController struct {
	UserService *models.UserService
	UserUrl     string
	oauth2.Config
}

func NewGoogleController(userService *models.UserService) (GoogleController, error) {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	if clientId == "" {
		return GoogleController{}, errors.New("google client key required")
	}

	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientSecret == "" {
		return GoogleController{}, errors.New("google secret key required")
	}

	return GoogleController{
		UserService: userService,
		UserUrl:     "https://www.googleapis.com/oauth2/v3/tokeninfo",
		Config: oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  "http://localhost:3000/auth/google/callback",
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"}},
	}, nil
}

func (controller *GoogleController) Signin(w http.ResponseWriter, r *http.Request) {
	authUrl := controller.AuthCodeURL("", oauth2.AccessTypeOffline)

	// TODO: Generate random string for state
	// TODO: Store it on cache ?

	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
}

func (controller *GoogleController) Callback(w http.ResponseWriter, r *http.Request) {
	token, err := controller.Exchange(context.TODO(), r.URL.Query().Get("code"))
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

func (controller *GoogleController) generateAccessTokenWithRefreshToken(refreshToken string) (*oauth2.Token, error) {
	requestBodyMap := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     controller.ClientID,
		"client_secret": controller.ClientSecret,
		"refresh_token": refreshToken,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest(
		"POST",
		controller.Endpoint.TokenURL,
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
