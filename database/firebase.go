package database

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type FirebaseConfig struct {
	Path string
}

func OpenFirebase(config FirebaseConfig) (*firebase.App, error) {
	return firebase.NewApp(context.TODO(), nil, option.WithCredentialsFile(config.Path))
}
