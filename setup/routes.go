package setup

import (
	"github.com/AkifhanIlgaz/vocab-builder/oauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(controllers *controllers, oauthHandlers *oauth.OAuthHandlers, middlewares *middlewares) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts

		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// TODO: Uid middleware

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", controllers.UsersController.Signup)
		r.Post("/signin", controllers.UsersController.Signin)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AccessTokenMiddleware.AccessToken)
			r.Post("/signout", controllers.UsersController.Signout)
		})
		r.Get("/refresh", controllers.UsersController.RefreshAccessToken)

		// TODO: Check if user is new on Signin endpoints
		r.Get("/google", controllers.GoogleController.Signin)
		r.Get("/google/callback", controllers.GoogleController.Callback)

	})

	r.Route("/box", func(r chi.Router) {
		r.Use(middlewares.AccessTokenMiddleware.AccessToken)
		r.Get("/today", controllers.BoxController.GetTodaysWords)
		r.Post("/levelup/{id}", controllers.BoxController.LevelUp)
		r.Post("/leveldown/{id}", controllers.BoxController.LevelDown)
	})

	return r
}
