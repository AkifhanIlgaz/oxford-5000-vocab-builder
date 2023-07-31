package database

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func OpenFirebase(credentialsFile string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(credentialsFile)
	return firebase.NewApp(context.TODO(), nil, opt)
}
