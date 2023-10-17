package setup

import ctrls "github.com/AkifhanIlgaz/vocab-builder/controllers"

type middlewares struct {
	AccessTokenMiddleware *ctrls.AccessTokenMiddleware
}

func Middlewares(services *services) *middlewares {
	return &middlewares{
		AccessTokenMiddleware: &ctrls.AccessTokenMiddleware{
			TokenService: services.TokenService,
		},
	}
}
