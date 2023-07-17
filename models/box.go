package models

import (
	"errors"
	"sort"
	"time"

	"github.com/boltdb/bolt"
)

/*
	Box
		Level 0 => Every day
		Level 1 => Every other day
		Level 2 => Week
		Level 3 => Month
		Level 4 => 3 Months
*/

var (
	ErrMaxLevel = errors.New("max level ")
	ErrMinLevel = errors.New("min level")
)

type BoxLevel int

const (
	BoxLevel0 BoxLevel = iota
	BoxLevel1
	BoxLevel2
	BoxLevel3
	BoxLevel4
)

const (
	Day   = 24 * time.Hour
	Week  = 7 * Day
	Month = 30 * Day //  Month is fixed to 30 days
)

const wordBoxLen = 5948

type WordBox []Word

func NewWordBox() WordBox {
	return make([]Word, wordBoxLen)
}

func (wb WordBox) getWordIds() []Word {
	var wordIds []Word

	for _, word := range wb {
		if word.NextRep.Before(time.Now()) {
			wordIds = append(wordIds, word)
		}
	}

	sort.Slice(wordIds, func(i, j int) bool {
		return wordIds[i].BoxLevel > wordIds[j].BoxLevel
	})

	return wordIds
}

type Word struct {
	Id          int
	BoxLevel    BoxLevel
	RepOnLevel1 int
	NextRep     time.Time
}

func (w *Word) boxLevelUp() error {
	boxLevel := w.BoxLevel

	if boxLevel == BoxLevel1 && w.RepOnLevel1 < 3 {
		w.RepOnLevel1++
		w.NextRep = w.NextRep.Add(2 * Day)
		return nil
	}

	// Level up
	if boxLevel < 4 {
		w.BoxLevel++
		switch w.BoxLevel {
		case BoxLevel1:
			w.NextRep = w.NextRep.Add(2 * Day)
		case BoxLevel2:
			w.NextRep = w.NextRep.Add(1 * Week)
		case BoxLevel3:
			w.NextRep = w.NextRep.Add(1 * Month)
		case BoxLevel4:
			w.NextRep = w.NextRep.Add(3 * Month)
		}
	} else {
		w.NextRep = w.NextRep.Add(3 * Month)
		return ErrMaxLevel
	}

	return nil
}

func (w *Word) boxLevelDown() error {
	boxLevel := w.BoxLevel
	if boxLevel <= BoxLevel0 {
		return ErrMinLevel
	}

	w.BoxLevel--
	switch w.BoxLevel {
	case BoxLevel0:
		w.NextRep = time.Now()
	case BoxLevel1:
		w.NextRep = w.NextRep.Add(2 * Day)
		w.RepOnLevel1 = 0
	case BoxLevel2:
		w.NextRep = w.NextRep.Add(1 * Week)
	case BoxLevel3:
		w.NextRep = w.NextRep.Add(1 * Month)
	}

	return nil
}

type BoxService struct {
	DB *bolt.DB
}

func (service *BoxService) GetTodaysWords(userId int) ([]*Word, error) {
	panic("Implement this function")
}

func (service *BoxService) GetWordsByLevel(userId int, boxLevel int) ([]*Word, error) {
	panic("Implement")
}
