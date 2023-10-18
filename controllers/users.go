package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/AkifhanIlgaz/vocab-builder/errors"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"golang.org/x/crypto/bcrypt"
)

type UsersController struct {
	WordService  *models.WordService
	BoxService   *models.BoxService
	UserService  *models.UserService
	TokenService *models.TokenService
}

type AuthenticationResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

/*
	TODO: When user is authenticated return
	{
		accessToken: "ey........"
		refreshToken: "sdilfksşdflk"
	}
*/

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
	// Extract email and password
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get user from db
	user, err := controller.UserService.GetByEmail(email)
	if err != nil {
		if errors.As(err, errors.ErrUserNotExist) {
			// ? Redirect to signup page
			http.Error(w, "User doesn't exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Wrong password", http.StatusInternalServerError)
		return
	}

	// Create new access token and refresh token for the user
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
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (controller *UsersController) Signout(w http.ResponseWriter, r *http.Request) {
	uid := context.Uid(r.Context())

	err := controller.TokenService.DeleteRefreshToken(uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// TODO: Refresh token is random string. Don't use JWT for refresh token
// TODO: Use only refresh token to generate new access token.
// TODO: After refresh token is consumed generate new refresh token and invalidate the previous one.
func (controller *UsersController) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract refresh token from request
	refreshToken := r.FormValue("refreshToken")
	if refreshToken == "" {
		http.Error(w, "Refresh token required", http.StatusBadRequest)
		return
	}

	newAccessToken, newRefreshToken, err := controller.TokenService.RefreshAccessToken(refreshToken)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&AuthenticationResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

}
