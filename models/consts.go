package models

import (
	"os"
	"time"
)

const (
	Database        = "VocabBuilder"
	UsersCollection = "Users"
	WordsColleciton = "Words"
)

var (
	ExpireDuration time.Duration = 1 * time.Hour
)

var (
	Secret = []byte(os.Getenv("TOKEN_SECRET"))
)
