package main

import (
	"context"
	"os"

	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	client, err := models.Open(os.Getenv("MONGODB_URI"))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// wordsCollection := client.Database("Vocab-Builder").Collection("Words")
	// wordService := models.WordService{
	// 	WordCollection: wordsCollection,
	// }

}
