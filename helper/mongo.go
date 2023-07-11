package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/AkifhanIlgaz/vocab-builder/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Without concurrency 40 minutes
With concurrency 2 minutes

	When I tried concurrency, I get TLS handshake timeout error
*/
func InsertToMongo(urlsFile string, wordsCollection *mongo.Collection) {
	file, err := os.Open(urlsFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var urls []struct {
		Url string `json:"url"`
	}

	dec := json.NewDecoder(file)
	dec.Decode(&urls)

	start := time.Now()

	for id, url := range urls {
		wordInfo, err := parser.ParseWord(url.Url)
		if err != nil {
			retry(id, url.Url, wordsCollection)
		}
		wordInfo.Id = id
		fmt.Println("Inserting", url.Url)
		wordsCollection.InsertOne(context.TODO(), wordInfo)
	}

	fmt.Println("5947 words is parsed and inserted to mongo in", time.Since(start))
}

func retry(id int, wordUrl string, wordsCollection *mongo.Collection) {
	wordInfo, err := parser.ParseWord(wordUrl)
	if err != nil {
		fmt.Println("Cannot parse", wordUrl)

	}
	wordInfo.Id = id
	fmt.Println("Inserting", wordUrl)
	wordsCollection.InsertOne(context.TODO(), wordInfo)
}

func AddBoxFieldToDocuments(wordsCollection *mongo.Collection) {

	wordsCollection.UpdateMany(context.TODO(), bson.D{}, bson.D{
		{
			Key: "$set", Value: bson.D{{
				Key: "box", Value: 1,
			}},
		},
	})
}
