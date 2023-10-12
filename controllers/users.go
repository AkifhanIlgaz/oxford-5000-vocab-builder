package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type UsersController struct {
	WordService *models.WordService
	BoxService  *models.BoxService
	UserService *models.UserService
}

func (controller *UsersController) Signup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := controller.UserService.Create(email, password)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	token, err := models.NewIdToken(user.Uid.Hex())
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	}{
		User:  *user,
		Token: token,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
