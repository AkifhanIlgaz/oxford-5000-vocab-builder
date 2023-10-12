package setup

import (
	"github.com/go-chi/chi/v5"
)

func Routes(controllers *controllers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(controllers.UserMiddleware.SetUser)

	r.Post("/auth/signup", controllers.UsersController.Signup)

	r.Route("/box", func(r chi.Router) {
		r.Use(controllers.UserMiddleware.RequireUser)
		r.Get("/today", controllers.BoxController.GetTodaysWords)
		r.Post("/levelup/{id}", controllers.BoxController.LevelUp)
		r.Post("/leveldown/{id}", controllers.BoxController.LevelDown)
	})

	return r
}
