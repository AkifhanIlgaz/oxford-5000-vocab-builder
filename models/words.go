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

// wordsPerDay will be stored on User table or as cookie ?
func (service *WordService) GetWordPackage(packageId int, wordsPerDay int) []*WordInfo {
	var words []*WordInfo

	for _, wordId := range service.getWordIds(packageId, wordsPerDay) {
		word, err := service.GetWord(wordId)
		if err != nil {
			return nil
		}
		words = append(words, word)
	}

	return words
}

// packageId will be 1-indexed
// [wordsPerDay * (packageId -1), wordsPerDay * packageId)
func (service *WordService) getWordIds(packageId int, wordsPerDay int) []int {
	wordIds := []int{}

	for i := wordsPerDay * (packageId - 1); i < wordsPerDay*(packageId); i++ {
		wordIds = append(wordIds, i)
	}

	return wordIds
}
