package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type UsersController struct {
	UserService    *models.UserService
	SessionService *models.SessionService
	WordService    *models.WordService
	BoxService     *models.BoxService
}

func (uc UsersController) SignUp(w http.ResponseWriter, r *http.Request) {
	// Parse form
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Used parsed information to create new user
	user, err := uc.UserService.Create(email, password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			// TODO: Return with appropriate status code and error message
			http.Error(w, models.ErrEmailTaken.Error(), http.StatusNotFound)
			return
		}
	}

	// Create session token
	session, err := uc.SessionService.Create(user.Id)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// set cookie
	setCookie(w, CookieSession, session.Token)

	fmt.Fprintln(w, "User successfully created")
}

func (uc UsersController) SignIn(w http.ResponseWriter, r *http.Request) {
	// Parse data from request
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Check if a given email and password matches within database
	user, err := uc.UserService.Authenticate(email, password)
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
	session, err := uc.SessionService.Create(user.Id)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)

	fmt.Fprint(w, "Logged in successfully")
}

func (uc UsersController) SignOut(w http.ResponseWriter, r *http.Request) {
	// Read session token from request
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Delete cookie
	if err := uc.SessionService.Delete(token); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, CookieSession)

	fmt.Fprint(w, "Logged out")
}

func (uc UsersController) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	fmt.Fprint(w, user)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Delete this line when go to production
		w.Header().Set("Access-Control-Allow-Origin", "*")

		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithUser(r.Context(), user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			// TODO: Redirect
			http.Error(w, "Please login", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
