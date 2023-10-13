package models

import (
	"os"
	"time"
)

const (
	Database               = "VocabBuilder"
	UsersCollection        = "Users"
	WordsCollection        = "Words"
	RefreshTokenCollection = "RefreshTokens"
)

var (
	ExpireDuration time.Duration = 1 * time.Hour
)

var (
	Secret = []byte(os.Getenv("TOKEN_SECRET"))
)
