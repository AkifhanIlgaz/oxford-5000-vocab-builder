package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	panic("")
}

func (service *WordService) BoxLevelDown(wordId int) error {
	panic("")
}
