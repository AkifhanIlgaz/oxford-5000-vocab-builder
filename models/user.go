package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/AkifhanIlgaz/vocab-builder/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Auth *auth.Client
}

type UserService struct {
	Collection *mongo.Collection
}

func NewUserService(client *mongo.Client) UserService {
	collection := client.Database(Database).Collection(UsersCollection)

	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	// TODO: Check error
	collection.Indexes().CreateOne(context.TODO(), indexModel)

	return UserService{
		Collection: collection,
	}
}

type User struct {
	Uid          primitive.ObjectID `json:"uid" bson:"_id"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	CreatedAt    time.Time          `json:"-" bson:"createdAt"`
}

func (service *UserService) Create(email, password string) (*User, error) {
	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	user := User{
		Uid:          primitive.NewObjectID(),
		Email:        strings.TrimSpace(email),
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}

	// Insert into DB
	_, err = service.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.ErrEmailTaken
		}
		return nil, err
	}

	return &user, nil
}
