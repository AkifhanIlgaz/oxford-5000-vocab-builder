package main

import (
	"context"
	"os"

	"github.com/AkifhanIlgaz/vocab-builder/models"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mongoDB, err := models.OpenMongo(os.Getenv("MONGODB_URI"))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = mongoDB.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

}
