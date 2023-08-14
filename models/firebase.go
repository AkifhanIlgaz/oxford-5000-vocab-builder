package models

import (
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

// TODO: Add auth as a struct field
type FirebaseService struct {
	App  *firebase.App
	Auth *auth.Client
}

// TODO: Implement helper functions for Firebase
func (service *FirebaseService) User(idToken string) {

}
