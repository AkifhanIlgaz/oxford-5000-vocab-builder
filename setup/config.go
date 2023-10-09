package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AkifhanIlgaz/vocab-builder/database"
	"github.com/joho/godotenv"
)

type config struct {
	Mongo    database.MongoConfig
	Bolt     database.BoltConfig // UserHomeDir + FileName
	Firebase database.FirebaseConfig
}

func Config() (*config, error) {
	var config config

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	config.Mongo = database.MongoConfig{
		ConnectionURI: os.Getenv("MONGODB_URI"),
	}
	if config.Mongo.ConnectionURI == "" {
		return nil, fmt.Errorf("no Mongo config provided")
	}

	boltFileName := os.Getenv("BOLT_FILENAME")
	if boltFileName == "" {
		return nil, fmt.Errorf("no BoltDB config provided")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}
	config.Bolt.Path = filepath.Join(home, boltFileName)

	config.Firebase.Path = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	return &config, nil
}
