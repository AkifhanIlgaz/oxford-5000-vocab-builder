package controllers

import (
	"errors"
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

func (u Users) SignUp(w http.ResponseWriter, r *http.Request) {
	// Parse form
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Used parsed information to create new user
	user, err := u.UserService.Create(email, password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			// TODO: Return with appropriate status code and error message
			http.Error(w, models.ErrEmailTaken.Error(), http.StatusNotFound)
			return
		}
	}

	// Create session token
	session, err := u.SessionService.Create(user.Id)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// set cookie
	setCookie(w, CookieSession, session.Token)

	fmt.Fprintln(w, "User successfully created")
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	// Parse data from request
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Check if a given email and password matches within database
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrWrongPassword) {
			// TODO: Appropriate status code
			http.Error(w, "Wrong password", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// create session token
	session, err := u.SessionService.Create(user.Id)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)

	fmt.Fprint(w, "Logged in successfully")
}
