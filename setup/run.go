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

	services, err := Services(databases)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	controllers, err := Controllers(services)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	r := Routes(controllers)

	fmt.Println("Starting server on", port)
	return http.ListenAndServe(port, r)
}
