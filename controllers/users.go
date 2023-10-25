package controllers

import (
	"encoding/json"
	"fmt"
	"io"
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

func (controller *UsersController) Signup(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	user, err := controller.UserService.Create(data.Email, data.Password)
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

	err = json.NewEncoder(w).Encode(&AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (controller *UsersController) Signin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	// Get user from db
	user, err := controller.UserService.GetByEmail(data.Email)
	if err != nil {
		if errors.As(err, errors.ErrUserNotExist) {
			// ? Redirect to signup page
			http.Redirect(w, r, "http://localhost:8100/signin", http.StatusUnauthorized)
			return
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password))
	if err != nil {
		http.Error(w, "Wrong password", http.StatusInternalServerError)
		return
	}

	// Create new access token and refresh token for the user
	accessToken, err := controller.TokenService.NewAccessToken(user.Uid.Hex())
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%v", err)
		return
	}

	refreshToken, err := controller.TokenService.NewRefreshToken(user.Uid.Hex())
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = json.NewEncoder(w).Encode(&AuthenticationResponse{
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

func (controller *UsersController) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refresh_token")
	if refreshToken == "" {
		http.Error(w, "Refresh token required", http.StatusBadRequest)
		return
	}

	fmt.Println(refreshToken)

	newAccessToken, newRefreshToken, err := controller.TokenService.RefreshAccessToken(refreshToken)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&AuthenticationResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

}
