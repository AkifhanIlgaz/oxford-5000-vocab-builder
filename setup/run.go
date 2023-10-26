package setup

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/oauth"
)

func Run(port string) error {
	config, err := Config()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	databases, err := Databases(config)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	oauthHandlers, err := oauth.Setup()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	services := Services(databases)

	controllers := Controllers(services)

	middlewares := Middlewares(services, oauthHandlers.Google)

	r := Routes(controllers, oauthHandlers, middlewares)

	fmt.Println("Starting server on", port)
	return http.ListenAndServe(port, r)
}
