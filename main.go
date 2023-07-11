package main

import (
	"context"
	"os"

	"github.com/AkifhanIlgaz/vocab-builder/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := database.Open(os.Getenv("MONGODB_URI"))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

}
