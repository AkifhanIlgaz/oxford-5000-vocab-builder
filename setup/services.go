package setup

import (
	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type services struct {
	WordService *models.WordService
	BoxService  *models.BoxService
	UserService *models.UserService
}

func Services(databases *databases) (*services, error) {

	wordService := models.WordService{
		Client: databases.Mongo,
	}

	userService := models.NewUserService(databases.Mongo)

	boxService := models.BoxService{
		DB: databases.Bolt,
	}

	return &services{
		WordService: &wordService,
		BoxService:  &boxService,
		UserService: &userService,
	}, nil
}
