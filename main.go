package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/AkifhanIlgaz/vocab-builder/controllers"
	"github.com/AkifhanIlgaz/vocab-builder/database"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/AkifhanIlgaz/vocab-builder/parser"
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
	Firebase database.FirebaseConfig
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

	cfg.Firebase.Path = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

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

func initServices(cfg config) (*mongo.Client, *sql.DB, *bolt.DB, *firebase.App, error) {
	mongo, err := database.OpenMongo(cfg.Mongo)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	fmt.Println("Connected to mongo")

	postgres, err := database.OpenPostgres(cfg.Postgres)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	fmt.Println("Connected to postgres")

	bolt, err := database.OpenBolt(cfg.Bolt)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	fmt.Println("Connected to bolt")

	firebase, err := database.OpenFirebase(cfg.Firebase)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return mongo, postgres, bolt, firebase, nil
}

func run(cfg config) error {
	mongo, postgres, bolt, firebaseApp, err := initServices(cfg)
	if err != nil {
		return err
	}

	// TODO: Delete this struct
	type FirebaseService struct {
		App *firebase.App
	}

	firebaseService := FirebaseService{
		App: firebaseApp,
	}
	fmt.Println(firebaseService)

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
	// r.Use(userMiddleware.SetUser)

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
		// TODO: Create new wordbox when user signs up
		r.Post("/new", boxController.NewWordBox)
		r.Get("/today", boxController.GetTodaysWords)
		r.Post("/levelup/{id}", boxController.LevelUp)
		r.Post("/leveldown/{id}", boxController.LevelDown)
	})

	r.Get("/words/{id}", wordsController.WordWithId)
	r.Get("/words/parse/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		word, _ := wordsController.WordService.GetWord(id)
		newParsed, _ := parser.ParseWord(word.Source)

		enc := json.NewEncoder(w)
		enc.Encode(newParsed)
	})
	// TODO: Get 10 random words
	r.Get("/words/random", func(w http.ResponseWriter, r *http.Request) {
		var words []*models.WordInfo
		source := rand.NewSource(time.Now().Unix())
		random := rand.New(source)

		for i := 0; i < 10; i++ {
			id := random.Intn(5948)
			word, _ := wordsController.WordService.GetWord(id)
			words = append(words, word)
		}

		enc := json.NewEncoder(w)
		enc.Encode(words)

	})

	fmt.Println("Starting server on", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
