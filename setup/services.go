package setup

import (
	"context"
	"fmt"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type services struct {
	AuthService *models.AuthService
	WordService *models.WordService
	BoxService  *models.BoxService
	UserService *models.UserService
}

func Services(databases *databases) (*services, error) {
	auth, err := databases.Firebase.Auth(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("setup services | auth: %w", err)
	}

	authService := models.AuthService{
		Auth: auth,
	}

	wordService := models.WordService{
		Client: databases.Mongo,
	}

	userService := models.NewUserService(databases.Mongo)

	boxService := models.BoxService{
		DB: databases.Bolt,
	}

	return &services{
		AuthService: &authService,
		WordService: &wordService,
		BoxService:  &boxService,
		UserService: &userService,
	}, nil
}
