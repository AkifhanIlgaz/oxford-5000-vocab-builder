package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	ctx "github.com/AkifhanIlgaz/vocab-builder/context"
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
	}

	err = json.NewEncoder(w).Encode((user))
	if err != nil {
		fmt.Println(err)
		return
	}
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
		uid := ctx.Uid(r.Context())
		if uid == "" {
			http.Error(w, "Please login", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
