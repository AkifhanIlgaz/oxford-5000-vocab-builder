package setup

import ctrls "github.com/AkifhanIlgaz/vocab-builder/controllers"

type controllers struct {
	BoxController    *ctrls.BoxController
	UsersController  *ctrls.UsersController
	GoogleController *ctrls.GoogleController
}

func Controllers(services *services) *controllers {
	boxController := ctrls.BoxController{
		BoxService:  services.BoxService,
		WordService: services.WordService,
	}

	usersController := ctrls.UsersController{
		WordService:  services.WordService,
		BoxService:   services.BoxService,
		UserService:  services.UserService,
		TokenService: services.TokenService,
	}

	// ! Check Error
	googleController, _ := ctrls.NewGoogleController(services.UserService)

	return &controllers{
		BoxController:    &boxController,
		UsersController:  &usersController,
		GoogleController: &googleController,
	}
}
