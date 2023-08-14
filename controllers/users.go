package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	ctx "github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type UsersController struct {
	UserService          *models.UserService
	SessionService       *models.SessionService
	WordService          *models.WordService
	BoxService           *models.BoxService
	EmailService         *models.EmailService
	PasswordResetService *models.PasswordResetService
}

// OK
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

	// http.Redirect(w, r, "/profile", http.StatusFound)
}

// OK
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

	http.Redirect(w, r, "/profile", http.StatusFound)
}

// OK
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

func (uc UsersController) Profile(w http.ResponseWriter, r *http.Request) {
	user := ctx.User(r.Context())

	fmt.Fprint(w, user)
}

func (uc UsersController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	passwordReset, err := uc.PasswordResetService.Create(data.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": {passwordReset.Token},
	}
	resetUrl := "localhost:8100" + "/reset-password?" + vals.Encode()

	err = uc.EmailService.ForgotPassword(data.Email, resetUrl)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusOK)
}

func (uc UsersController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}

	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := uc.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		// TODO: Handle different kind of errors
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = uc.UserService.UpdatePassword(user.Id, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := uc.SessionService.Create(user.Id)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/profile", http.StatusOK)
}

type UserMiddleware struct {
	AuthService *models.AuthService
}

func readToken(r *http.Request) string {
	return r.URL.Query().Get("token")
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := readToken(r)

		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		authToken, err := umw.AuthService.Auth.VerifyIDToken(context.TODO(), token)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}
		newContext := ctx.WithUid(r.Context(), authToken.UID)
		r = r.WithContext(newContext)

		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ctx.User(r.Context())
		if user == nil {
			// TODO: Redirect
			http.Error(w, "Please login", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
