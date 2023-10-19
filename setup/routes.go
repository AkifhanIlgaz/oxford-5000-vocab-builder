package setup

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

func Routes(controllers *controllers, middlewares *middlewares) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", controllers.UsersController.Signup)
		r.Post("/signin", controllers.UsersController.Signin)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AccessTokenMiddleware.AccessToken)
			r.Post("/signout", controllers.UsersController.Signout)
		})

		r.Get("/refresh", controllers.UsersController.RefreshAccessToken)
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

	oauth2.NewClient(nil, oauth2.StaticTokenSource(&oauth2.Token{}))

	return r
}
