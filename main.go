package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"github.com/AkifhanIlgaz/vocab-builder/controllers"
	"github.com/AkifhanIlgaz/vocab-builder/database"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/boltdb/bolt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type config struct {
	Mongo    database.MongoConfig
	Bolt     database.BoltConfig // UserHomeDir + FileName
	Firebase database.FirebaseConfig
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

	cfg.Firebase.Path = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

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

func initServices(cfg config) (*mongo.Client, *bolt.DB, *firebase.App, error) {
	mongo, err := database.OpenMongo(cfg.Mongo)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println("Connected to mongo")

	bolt, err := database.OpenBolt(cfg.Bolt)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Println("Connected to bolt")

	app, err := database.OpenFirebase(cfg.Firebase)
	if err != nil {
		return nil, nil, nil, err
	}

	return mongo, bolt, app, nil
}

func run(cfg config) error {
	mongo, bolt, app, err := initServices(cfg)
	if err != nil {
		return err
	}

	auth, err := app.Auth(context.TODO())
	if err != nil {
		return err
	}

	authService := models.AuthService{
		Auth: auth,
	}

	wordService := models.WordService{
		Client: mongo,
	}

	boxService := models.BoxService{
		DB: bolt,
	}

	r := chi.NewRouter()

	boxController := controllers.BoxController{
		BoxService:  &boxService,
		WordService: &wordService,
	}

	userMiddleware := controllers.UserMiddleware{
		AuthService: &authService,
	}

	r.Use(userMiddleware.SetUser)

	r.Route("/box", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/today", boxController.GetTodaysWords)
		r.Post("/levelup/{id}", boxController.LevelUp)
		r.Post("/leveldown/{id}", boxController.LevelDown)
	})

	fmt.Println("Starting server on", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
