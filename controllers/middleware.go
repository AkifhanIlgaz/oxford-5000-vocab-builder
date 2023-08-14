package controllers

import (
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

// r.Get("/idtoken", func(w http.ResponseWriter, r *http.Request) {
// 		auth, err := firebaseService.App.Auth(context.TODO())
// 		if err != nil {
// 			panic(err)
// 		}

// 		token, err := auth.VerifyIDToken(context.TODO(), r.URL.Query().Get("token"))
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(token.UID)

// 	})

type FirebaseMiddleware struct {
	App *models.FirebaseService
}

func (fmw *FirebaseMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, err := fmw.App.App.Auth(context.TODO())
		if err != nil {
			panic(err)
		}

		token, err := readCookie(r, TokenSession)
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

func SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := readCookie(r, TokenSession)
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
