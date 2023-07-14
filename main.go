package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AkifhanIlgaz/vocab-builder/controllers"
	"github.com/AkifhanIlgaz/vocab-builder/database"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/boltdb/bolt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
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
	config, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	err = run(config)
	if err != nil {
		panic(err)
	}
}

func initServices(cfg config) (*mongo.Client, *sql.DB, *bolt.DB, error) {
	// Check your allowed IP address for mongo
	mongo, err := database.OpenMongo(cfg.Mongo)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println("Connected to mongo")

	postgres, err := database.OpenPostgres(cfg.Postgres)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println("Connected to postgres")

	bolt, err := database.OpenBolt(cfg.Bolt)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println("Connected to bolt")

	return mongo, postgres, bolt, nil
}

func run(cfg config) error {
	mongo, postgres, bolt, err := initServices(cfg)
	if err != nil {
		return err
	}

	userService := models.UserService{
		DB: postgres,
	}

	wordService := models.WordService{
		Client: mongo,
	}

	boxService := models.BoxService{
		DB: bolt,
	}

	sessionService := models.SessionService{
		DB: postgres,
	}

	r := chi.NewRouter()

	usersController := controllers.Users{
		UserService:    &userService,
		WordService:    &wordService,
		BoxService:     &boxService,
		SessionService: &sessionService,
	}

	r.Post("/signup", usersController.SignUp)
	r.Post("/signin", usersController.SignIn)
	r.Post("/signout", usersController.SignOut)

	fmt.Println("Starting server on", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
