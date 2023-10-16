package setup

import (
	"time"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type services struct {
	WordService  *models.WordService
	BoxService   *models.BoxService
	UserService  *models.UserService
	TokenService *models.TokenService
}

func Services(databases *databases) (*services, error) {

	wordService := models.NewWordService(databases.Mongo)

	userService := models.NewUserService(databases.Mongo)

	tokenService := models.NewTokenService(databases.Mongo, 15*time.Minute)

	boxService := models.BoxService{
		DB: databases.Bolt,
	}

	return &services{
		WordService:  &wordService,
		BoxService:   &boxService,
		UserService:  &userService,
		TokenService: &tokenService,
	}, nil
}
