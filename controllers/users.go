package controllers

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type Users struct {
	UserService    *models.UserService
	SessionService *models.SessionService
	WordService    *models.WordService
	BoxService     *models.BoxService
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// Parse form
	// Used parsed information to create new user
	// Create session token
	// set cookie

	fmt.Println(r.FormValue("email"), r.FormValue("password"))
	fmt.Fprintln(w, "Create user endpoint")
}
