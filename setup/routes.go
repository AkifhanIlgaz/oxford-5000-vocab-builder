package setup

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/AkifhanIlgaz/vocab-builder/oauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(controllers *controllers, oauthHandlers *oauth.OAuthHandlers, middlewares *middlewares) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "idToken", "fileExtension", "type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", controllers.UsersController.Signup)
		r.Post("/signin", controllers.UsersController.Signin)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AccessTokenMiddleware.AccessToken)
			r.Post("/signout", controllers.UsersController.Signout)
		})
		r.Get("/refresh", controllers.UsersController.RefreshAccessToken)

		r.Route("/signin", func(r chi.Router) {
			// TODO: Check if user is new on Signin endpoints
			r.Get("/github", oauthHandlers.Github.Signin)
			r.Get("/github/callback", oauthHandlers.Github.Callback)
			// TODO: Add refresh endpoint

			r.Get("/google", oauthHandlers.Google.Signin)
			r.Get("/google/callback", oauthHandlers.Google.Callback)
			r.Get("/google/refresh", oauthHandlers.Google.GenerateAccessTokenWithRefreshToken)
		})

	})

	r.Route("/test", func(r chi.Router) {

		r.Use(middlewares.AccessTokenMiddleware.AccessToken)
		r.Get("/access-token-middleware", func(w http.ResponseWriter, r *http.Request) {
			uid := context.Uid(r.Context())

			fmt.Fprint(w, uid)
		})
	})

	r.Get("/x", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})
	r.Route("/box", func(r chi.Router) {
		r.Get("/today", controllers.BoxController.GetTodaysWords)
		r.Post("/levelup/{id}", controllers.BoxController.LevelUp)
		r.Post("/leveldown/{id}", controllers.BoxController.LevelDown)
	})

	return r
}

func Auth(r *chi.Mux, controllers *controllers, oauthHandlers *oauth.OAuthHandlers, middlewares *middlewares) error {

	return nil
}
