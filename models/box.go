package models

import (
	"errors"
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

type Word struct {
	Id          int
	BoxLevel    BoxLevel
	RepOnLevel1 int
	NextRep     time.Time
}

func (w *Word) BoxLevelUp() error {
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
		return errors.New("max level reached")
	}

	// TODO: Create constants for box levels and additional times for each box
	// TODO: Box level [0,4]
	// TODO: Word must be repeated 3 times on level 1 before moving to level 2.
	// Create additional field to determine how many times the word is repeated in level 2.
	// Reset this field when word is level down

	return nil
}

func (w *Word) BoxLevelDown() error {
	return nil
}

type BoxService struct {
	DB *bolt.DB
}
