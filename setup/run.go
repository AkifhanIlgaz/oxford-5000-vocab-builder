package setup

import (
	"fmt"
	"net/http"
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

	services := Services(databases)

	controllers := Controllers(services)

	middlewares := Middlewares(services)

	r := Routes(controllers, middlewares)

	fmt.Println("Starting server on", port)
	return http.ListenAndServe(port, r)
}
