package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AkifhanIlgaz/vocab-builder/database"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

// TODO: Add SMTP and CSRF config
type config struct {
	Postgres database.PostgresConfig
	Mongo    database.MongoConfig
	Bolt     database.BoltConfig // UserHomeDir + FileName
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

	cfg.Postgres = database.PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DBNAME"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
	if cfg.Postgres.Host == "" && cfg.Postgres.Port == "" {
		return cfg, fmt.Errorf("no Postgres config provided")
	}

	cfg.Mongo = database.MongoConfig{
		ConnectionURI: os.Getenv("MONGODB_URI"),
	}
	if cfg.Mongo.ConnectionURI == "" {
		return cfg, fmt.Errorf("no Mongo config provided")
	}

	boltFileName := os.Getenv("BOLT_FILENAME")
	if boltFileName == "" {
		return cfg, fmt.Errorf("no BoltDB config provided")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, fmt.Errorf("load env: %w", err)
	}
	cfg.Bolt.Path = filepath.Join(home, boltFileName)

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func main() {

}
