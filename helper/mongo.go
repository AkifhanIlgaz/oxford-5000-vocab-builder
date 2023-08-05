package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/AkifhanIlgaz/vocab-builder/parser"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Without concurrency 40 minutes
With concurrency 2 minutes

	When I tried concurrency, I get TLS handshake timeout error
*/

func UpdatePhonetics(collection *mongo.Collection) {
	words := []models.WordInfo{}

	cur, err := collection.Find(context.TODO(), bson.M{
		"header.partofspeech": "verb",
	})
	if err != nil {
		panic(err)
	}

	err = cur.All(context.TODO(), &words)
	if err != nil {
		panic(err)
	}
	for _, word := range words {
		parsedWord, _ := parser.ParseWord(word.Source)
		fmt.Println(parsedWord.Word)
		collection.UpdateOne(context.TODO(), bson.M{
			"source": word.Source,
		}, bson.D{
			{"$set", bson.D{
				{"header.audio", parsedWord.Header.Audio},
			}},
		})

	}
}

func EdgeCaseSenseSingle(collection *mongo.Collection) {

	words := []models.WordInfo{}

	cur, _ := collection.Find(context.TODO(), bson.D{
		{
			"definitions", nil,
		},
	})

	cur.All(context.TODO(), &words)

	for _, word := range words {
		retry(word.Id, word.Source, collection)
	}
	fmt.Println(words)
}

func DeleteNullDefinitions(collection *mongo.Collection) {

	res, err := collection.DeleteMany(context.TODO(), bson.D{
		{
			"definitions", nil,
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.DeletedCount)
}

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
		wordInfo.Source = url.Url
		fmt.Println("Inserting", url)
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

func RenameBoxField(wordsCollection *mongo.Collection) {
	wordsCollection.UpdateMany(context.TODO(), bson.D{}, bson.D{
		{
			Key: "$rename", Value: bson.D{{
				Key: "box", Value: "boxLevel",
			}},
		},
	})
}

func WithConcurrency() {
	godotenv.Load()

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	file, err := os.Open("word_database/urls.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var urls []struct {
		id  int
		Url string `json:"url"`
	}

	wordCollection := client.Database("VocabBuilder").Collection("Words")

	dec := json.NewDecoder(file)
	dec.Decode(&urls)

	start := time.Now()

	var wg sync.WaitGroup

	total := len(urls)
	ok := 0
	false := 0

	errorss := []struct {
		id  int
		Url string `json:"url"`
	}{}

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url struct {
			id  int
			Url string `json:"url"`
		}) {
			defer wg.Done()
			word, err := parser.ParseWord(url.Url)
			word.Id = i
			word.Source = url.Url
			if err != nil {
				fmt.Println("error", url.Url)
				url.id = i
				errorss = append(errorss, url)
			}
			_, err = wordCollection.InsertOne(context.TODO(), word)
			if err != nil {
				fmt.Println("error", url.Url)
				url.id = i
				errorss = append(errorss, url)
			}
			fmt.Println(url.Url)
			ok++
		}(i, url)

	}

	wg.Wait()

	for _, url := range errorss {

		word, err := parser.ParseWord(url.Url)
		word.Id = url.id
		if err != nil {
			fmt.Println(url.Url)
		}
		_, err = wordCollection.InsertOne(context.TODO(), word)

		if err != nil {
			fmt.Println(url.Url)
		}
		fmt.Println(url.Url)
		ok++
	}

	fmt.Println("total", total)
	fmt.Println("ok", ok)
	fmt.Println("false", false)
	fmt.Println(time.Since(start))
}
