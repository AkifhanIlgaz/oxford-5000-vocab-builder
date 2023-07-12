package main

import (
	"fmt"

	"github.com/AkifhanIlgaz/vocab-builder/database"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

// TODO: Add SMTP and CSRF config
type config struct {
	Postgres database.PostgresConfig
	Mongo    database.MongoConfig
	Bolt     database.BoltConfig // Use home dir
	Server   struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, fmt.Errorf("load env: %w", err)
	}

	return config{}, nil
}

func main() {

}
