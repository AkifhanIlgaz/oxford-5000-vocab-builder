package main

import (
	"context"
	"os"

	"github.com/AkifhanIlgaz/vocab-builder/helper"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/joho/godotenv"
)

const (
	withIdioms    = "https://www.oxfordlearnersdictionaries.com/definition/english/about_2"
	withoutIdioms = "https://www.oxfordlearnersdictionaries.com/definition/english/across_2"
	diff          = "https://www.oxfordlearnersdictionaries.com/definition/english/reject_1"
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

	wordsCollection := client.Database("Vocab-Builder").Collection("Words")
	helper.InsertToMongo("./word_database/urls.json", wordsCollection)

}
