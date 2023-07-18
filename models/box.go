package models

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
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

const (
	Day   = 24 * time.Hour
	Week  = 7 * Day
	Month = 30 * Day //  Month is fixed to 30 days
)

const wordBoxLen = 5948
const boxBucket = "Boxes"

type WordBox []Word

type Word struct {
	Id          int
	BoxLevel    int
	RepOnLevel1 int
	NextRep     time.Time
}

func (w *Word) levelUp() error {
	boxLevel := w.BoxLevel

	if boxLevel == 1 && w.RepOnLevel1 < 2 {
		w.RepOnLevel1 += 1
		w.NextRep = w.NextRep.Add(2 * Day)
		return nil
	}

	// Level up
	if boxLevel < 4 {
		w.BoxLevel++
		switch w.BoxLevel {
		case 1:
			w.NextRep = w.NextRep.Add(2 * Day)
		case 2:
			w.NextRep = w.NextRep.Add(1 * Week)
		case 3:
			w.NextRep = w.NextRep.Add(1 * Month)
		case 4:
			w.NextRep = w.NextRep.Add(3 * Month)
		}
	} else {
		w.NextRep = w.NextRep.Add(3 * Month)
		return ErrMaxLevel
	}

	return nil
}

func (w *Word) levelDown() error {
	boxLevel := w.BoxLevel
	if boxLevel <= 0 {
		return ErrMinLevel
	}

	w.BoxLevel--
	switch w.BoxLevel {
	case 0:
		w.NextRep = time.Now()
	case 1:
		w.NextRep = w.NextRep.Add(2 * Day)
		w.RepOnLevel1 = 0
	case 2:
		w.NextRep = w.NextRep.Add(1 * Week)
	case 3:
		w.NextRep = w.NextRep.Add(1 * Month)
	}

	return nil
}

type BoxService struct {
	DB *bolt.DB
}

func (service *BoxService) CreateWordBox(userId int) error {
	err := service.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boxBucket))
		var wordBox WordBox
		for i := 0; i < wordBoxLen; i++ {
			wordBox = append(wordBox, Word{
				Id:      i,
				NextRep: time.Now(),
			})
		}

		return b.Put(itob(userId), serializeWordBox(wordBox))
	})

	if err != nil {
		return fmt.Errorf("create wordbox: %w", err)
	}

	return nil
}

func (service *BoxService) getWordBox(userId int) (WordBox, error) {
	var wordBox WordBox

	err := service.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boxBucket))
		wordBox = deserializeWordBox(b.Get(itob(userId)))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("get wordBox: %w", err)
	}

	return wordBox, nil
}

func (service *BoxService) GetTodaysWords(userId int) ([]Word, error) {
	var todaysWords []Word

	wordBox, err := service.getWordBox(userId)
	if err != nil {
		return nil, fmt.Errorf("get todays words: %w", err)
	}

	for _, word := range wordBox {
		if word.NextRep.Before(time.Now()) {
			todaysWords = append(todaysWords, word)
		}
	}

	sort.Slice(todaysWords, func(i, j int) bool {
		return todaysWords[i].BoxLevel > todaysWords[j].BoxLevel
	})

	return todaysWords, nil
}

func (service *BoxService) updateWordBox(userId int, wordBox WordBox) error {
	err := service.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boxBucket))
		return b.Put(itob(userId), serializeWordBox(wordBox))
	})

	if err != nil {
		return fmt.Errorf("update wordbox: %w", err)
	}

	return nil
}

func (service *BoxService) LevelUp(userId int, wordId int) error {
	wordBox, err := service.getWordBox(userId)
	if err != nil {
		return fmt.Errorf("level up: %w", err)
	}
	w := &wordBox[wordId]
	w.levelUp()

	return service.updateWordBox(userId, wordBox)
}

func (service *BoxService) LevelDown(userId int, wordId int) error {
	wordBox, err := service.getWordBox(userId)
	if err != nil {
		return fmt.Errorf("level up: %w", err)
	}

	w := &wordBox[wordId]
	w.levelDown()
	return service.updateWordBox(userId, wordBox)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func serializeWordBox(wordBox WordBox) []byte {
	b, _ := json.MarshalIndent(wordBox, "", "  ")
	return b
}

func deserializeWordBox(b []byte) WordBox {
	var wb WordBox
	json.Unmarshal(b, &wb)
	return wb
}
