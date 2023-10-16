package models

import (
	"os"
)

const (
	Database               = "VocabBuilder"
	UsersCollection        = "Users"
	WordsCollection        = "Words"
	RefreshTokenCollection = "RefreshTokens"
)

var (
	Secret = []byte(os.Getenv("TOKEN_SECRET"))
)
