package setup

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/go-chi/chi/v5"
)

// TODO: Uid middleware
// TODO: Create /auth route for authentication purposes
// TODO: Parse Authorization Bearer Header

func Routes(controllers *controllers, middlewares *middlewares) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		// TODO: Implement sign in endpoint
		// TODO: Implement sign out endpoint
		r.Post("/signup", controllers.UsersController.Signup)

		// TODO: Implement /refresh endpoint && It will be used to generate new access token with refresh token
		r.Get("/refresh", func(w http.ResponseWriter, r *http.Request) {
			panic("// ! Implement")
		})

	})

	r.Route("/test", func(r chi.Router) {
		r.Use(middlewares.AccessTokenMiddleware.AccessToken)
		r.Get("/access-token-middleware", func(w http.ResponseWriter, r *http.Request) {
			uid := context.Uid(r.Context())

			fmt.Fprint(w, uid)
		})
	})

	r.Get("/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		// parse refresh token
		// token, err := (r.FormValue("refreshToken"))
		// if err != nil {
		// 	fmt.Fprint(w, "parse refresh token", err)
		// 	return
		// }

		// claims, ok := token.Claims.(*models.RefreshClaims)
		// if !ok {
		// 	fmt.Fprint(w, "parse refresh token interface")
		// 	return
		// }

		// idToken, err := controllers.UsersController.TokenService.RefreshAccessToken(claims.Subject, claims.RefreshToken)
		// if err != nil {
		// 	fmt.Fprint(w, "parse refresh token refresh")
		// 	return
		// }

		// fmt.Fprint(w, idToken)
	})

	r.Route("/box", func(r chi.Router) {
		r.Get("/today", controllers.BoxController.GetTodaysWords)
		r.Post("/levelup/{id}", controllers.BoxController.LevelUp)
		r.Post("/leveldown/{id}", controllers.BoxController.LevelDown)
	})

	return r
}
