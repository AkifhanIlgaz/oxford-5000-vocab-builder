package setup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/go-chi/chi/v5"
)

// TODO: Uid middleware

func Routes(controllers *controllers) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/auth/signup", controllers.UsersController.Signup)

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bearerToken := strings.Fields(r.Header.Get("Authorization"))
				if len(bearerToken) < 2 {
					fmt.Fprint(w, "invalid bearer token")
				}

				fmt.Println(bearerToken[1])

				token, err := models.ParseIdToken(bearerToken[1])
				if err != nil {
					fmt.Fprint(w, "parse id token")
					return
				}

				enc := json.NewEncoder(w)
				err = enc.Encode(token.Claims)
				if err != nil {
					fmt.Println(err)
				}
			})
		})
		r.Get("/auth/idtoken", nil)
	})

	r.Route("/box", func(r chi.Router) {
		r.Get("/today", controllers.BoxController.GetTodaysWords)
		r.Post("/levelup/{id}", controllers.BoxController.LevelUp)
		r.Post("/leveldown/{id}", controllers.BoxController.LevelDown)
	})

	return r
}
