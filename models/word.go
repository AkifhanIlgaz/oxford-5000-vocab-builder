package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WordInfo struct {
	Id          int          `json:"id"`
	Box         int          `json:"box"`
	Source      string       `json:"source"`
	Word        string       `json:"word"`
	Header      Header       `json:"header"`
	Definitions []Definition `json:"definitions"`
	Idioms      []Idiom      `json:"idioms"`
}

type Header struct {
	Audio struct {
		UK string `json:"UK"`
		US string `json:"US"`
	} `json:"audio"`
	PartOfSpeech string `json:"partOfSpeech"`
	CEFRLevel    string `json:"CEFRLevel"`
}

type Definition struct {
	Meaning  string   `json:"meaning"`
	Examples []string `json:"examples"`
}

type Idiom struct {
	Usage       string       `json:"usage"`
	Definitions []Definition `json:"definition"`
}

type WordService struct {
	Collection *mongo.Collection
}

func NewWordService(client *mongo.Client) WordService {
	return WordService{
		Collection: getCollection(client, WordsCollection),
	}
}

func (service *WordService) GetWord(id int) (*WordInfo, error) {
	var wordInfo WordInfo

	filter := bson.D{{Key: "id", Value: id}}
	err := service.Collection.FindOne(context.TODO(), filter).Decode(&wordInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("get word: %w", err)
		}
		return nil, fmt.Errorf("decoding word: %w", err)
	}

	return &wordInfo, nil
}

func (service *WordService) GetWordWithCollection(collection *mongo.Collection, id int) (*WordInfo, error) {
	// TODO: Pass the collection as parameter
	panic("")
}
