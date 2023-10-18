package setup

import (
	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type services struct {
	WordService  *models.WordService
	BoxService   *models.BoxService
	UserService  *models.UserService
	TokenService *models.TokenService
}

func Services(databases *databases) *services {

	wordService := models.NewWordService(databases.Mongo)

	userService := models.NewUserService(databases.Mongo)

	tokenService := models.NewTokenService(databases.Mongo, nil)

	boxService := models.BoxService{
		DB: databases.Bolt,
	}

	return &services{
		WordService:  &wordService,
		BoxService:   &boxService,
		UserService:  &userService,
		TokenService: &tokenService,
	}
}
