package setup

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/AkifhanIlgaz/vocab-builder/oauth"
	"github.com/go-chi/chi/v5"
)

func Routes(controllers *controllers, oauthHandlers *oauth.OAuthHandlers, middlewares *middlewares) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", controllers.UsersController.Signup)
		r.Post("/signin", controllers.UsersController.Signin)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AccessTokenMiddleware.AccessToken)
			r.Post("/signout", controllers.UsersController.Signout)
		})
		r.Get("/refresh", controllers.UsersController.RefreshAccessToken)

		r.Route("/signin", func(r chi.Router) {
			r.Get("/github", func(w http.ResponseWriter, r *http.Request) {})
			r.Get("/github/callback", func(w http.ResponseWriter, r *http.Request) {})

			r.Get("/google", oauthHandlers.Google.Signin)
			r.Get("/google/callback", oauthHandlers.Google.Callback)
			r.Get("/google/verify", oauthHandlers.Google.VerifyAccessToken)
		})

	})

	r.Route("/test", func(r chi.Router) {
		r.Use(middlewares.AccessTokenMiddleware.AccessToken)
		r.Get("/access-token-middleware", func(w http.ResponseWriter, r *http.Request) {
			uid := context.Uid(r.Context())

			fmt.Fprint(w, uid)
		})
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
