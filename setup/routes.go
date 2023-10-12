package setup

import (
	"github.com/go-chi/chi/v5"
)

// TODO: Uid middleware

func Routes(controllers *controllers) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/auth/signup", controllers.UsersController.Signup)

	r.Route("/box", func(r chi.Router) {
		r.Get("/today", controllers.BoxController.GetTodaysWords)
		r.Post("/levelup/{id}", controllers.BoxController.LevelUp)
		r.Post("/leveldown/{id}", controllers.BoxController.LevelDown)
	})

	return r
}
