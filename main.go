package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
	Postgres database.PostgresConfig
	Mongo    database.MongoConfig
	Bolt     database.BoltConfig // UserHomeDir + FileName
	SMTP     models.SMTPConfig
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

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.UserName = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

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

	firebase, err := database.OpenFirebase()
	fmt.Println("Connected to bolt")

	return mongo, bolt, firebase, nil
}

func run(cfg config) error {
	mongo, bolt, firebase, err := initServices(cfg)
	if err != nil {
		return err
	}

	userService := models.UserService{
		DB: firebase,
	}

	wordService := models.WordService{
		Client: mongo,
	}

	boxService := models.BoxService{
		DB: bolt,
	}

	sessionService := models.SessionService{
		DB: firebase,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	passwordResetService := models.PasswordResetService{
		DB: postgres,
	}

	r := chi.NewRouter()

	usersController := controllers.UsersController{
		UserService:          &userService,
		WordService:          &wordService,
		BoxService:           &boxService,
		SessionService:       &sessionService,
		EmailService:         &emailService,
		PasswordResetService: &passwordResetService,
	}

	wordsController := controllers.WordsController{
		WordService: &wordService,
	}

	boxController := controllers.BoxController{
		BoxService:  &boxService,
		WordService: &wordService,
	}

	userMiddleware := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	// All endpoints are working correctly
	r.Use(userMiddleware.SetUser)

	r.Post("/signup", usersController.SignUp)
	r.Post("/signin", usersController.SignIn)
	r.Post("/signout", usersController.SignOut)
	r.Post("/forgot-password", usersController.ForgotPassword)
	r.Post("/reset-password", usersController.ResetPassword)
	r.Route("/profile", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", usersController.Profile)
	})

	r.Route("/box", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		// TODO: Delete get wordbox endpoint
		r.Get("/", boxController.GetWordBox)
		r.Post("/new", boxController.NewWordBox)
		r.Get("/today", boxController.GetTodaysWords)
		r.Post("/levelup/{id}", boxController.LevelUp)
		r.Post("/leveldown/{id}", boxController.LevelDown)
	})

	r.Get("/words/{id}", wordsController.WordWithId)

	fmt.Println("Starting server on", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
