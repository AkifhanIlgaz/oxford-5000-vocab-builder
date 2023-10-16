package setup

import ctrls "github.com/AkifhanIlgaz/vocab-builder/controllers"

type controllers struct {
	BoxController   *ctrls.BoxController
	UsersController *ctrls.UsersController
}

func Controllers(services *services) (*controllers, error) {
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

	return &controllers{
		BoxController:   &boxController,
		UsersController: &usersController,
	}, nil
}
