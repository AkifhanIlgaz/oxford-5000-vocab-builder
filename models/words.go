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
	WordCollection *mongo.Collection
}

func (service *WordService) GetWord(id int) (*WordInfo, error) {
	var wordInfo WordInfo

	filter := bson.D{{Key: "id", Value: id}}
	err := service.WordCollection.FindOne(context.TODO(), filter).Decode(&wordInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("get word: %w", err)
		}
		return nil, fmt.Errorf("decoding word: %w", err)
	}

	return &wordInfo, nil
}

func (service *WordService) BoxLevelUp(wordId int) error {
	filter := bson.D{
		{
			Key:   "id",
			Value: wordId,
		},
	}

	update := bson.D{
		{
			Key: "$inc",
			Value: bson.D{{
				Key:   "boxLevel",
				Value: 1,
			}},
		},
	}

	result, err := service.WordCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("box level up: %w", err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("cannot found word with id: %v", wordId)
	}

	return nil
}

func (service *WordService) BoxLevelDown(wordId int) error {
	filter := bson.D{
		{
			Key:   "id",
			Value: wordId,
		},
	}

	update := bson.D{
		{
			Key: "$inc",
			Value: bson.D{{
				Key:   "boxLevel",
				Value: -1,
			}},
		},
	}

	result, err := service.WordCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("box level up: %w", err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("cannot found word with id: %v", wordId)
	}

	return nil
}
