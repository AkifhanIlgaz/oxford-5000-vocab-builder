package setup

import (
	ctrls "github.com/AkifhanIlgaz/vocab-builder/controllers"
	"github.com/AkifhanIlgaz/vocab-builder/oauth"
)

type middlewares struct {
	AccessTokenMiddleware *ctrls.AccessTokenMiddleware
}

func Middlewares(services *services, googleOauth *oauth.GoogleOAuth) *middlewares {
	return &middlewares{
		AccessTokenMiddleware: &ctrls.AccessTokenMiddleware{
			TokenService: services.TokenService,
			GoogleOauth:  googleOauth,
		},
	}
}
