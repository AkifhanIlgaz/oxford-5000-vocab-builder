package setup

import ctrls "github.com/AkifhanIlgaz/vocab-builder/controllers"

type controllers struct {
	BoxController   *ctrls.BoxController
	UserMiddleware  *ctrls.UserMiddleware
	UsersController *ctrls.UsersController
	// TODO: Add usercontroller
}

func Controllers(services *services) (*controllers, error) {
	boxController := ctrls.BoxController{
		BoxService:  services.BoxService,
		WordService: services.WordService,
	}

	usersController := ctrls.UsersController{
		WordService: services.WordService,
		BoxService:  services.BoxService,
		UserService: services.UserService,
	}

	userMiddleware := ctrls.UserMiddleware{
		AuthService: services.AuthService,
	}

	return &controllers{
		BoxController:   &boxController,
		UserMiddleware:  &userMiddleware,
		UsersController: &usersController,
	}, nil
}
