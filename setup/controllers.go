package setup

import ctrls "github.com/AkifhanIlgaz/vocab-builder/controllers"

type controllers struct {
	BoxController  *ctrls.BoxController
	UserMiddleware *ctrls.UserMiddleware
}

func Controllers(services *services) (*controllers, error) {
	boxController := ctrls.BoxController{
		BoxService:  services.BoxService,
		WordService: services.WordService,
	}

	userMiddleware := ctrls.UserMiddleware{
		AuthService: services.AuthService,
	}

	return &controllers{
		BoxController:  &boxController,
		UserMiddleware: &userMiddleware,
	}, nil
}
