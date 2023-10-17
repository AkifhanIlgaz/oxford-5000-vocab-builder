package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type UsersController struct {
	WordService  *models.WordService
	BoxService   *models.BoxService
	UserService  *models.UserService
	TokenService *models.TokenService
}

func (controller *UsersController) Signup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := controller.UserService.Create(email, password)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	accessToken, err := controller.TokenService.NewAccessToken(user.Uid.Hex())
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	refreshToken, err := controller.TokenService.NewRefreshToken(user.Uid.Hex())
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		User         models.User `json:"user"`
		AccessToken  string      `json:"accessToken"`
		RefreshToken string      `json:"refreshToken"`
	}{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (controller *UsersController) Signin(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract email and password

	// TODO: Check if user exists and credentials are true

	// TODO: Create new access token and refresh token for the user

	// TODO: Insert refresh token to DB
}

func (controller *UsersController) Signout(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse access token

	// TODO: Set Authorization header empty ? Can I set client headers from backend (like cookie) ?

	// TODO: Delete refresh token from DB
}

func (controller *UsersController) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract refresh token from request

	// TODO: Query database to check if refresh token credentials are true

	// TODO: If then, generate access token
}
