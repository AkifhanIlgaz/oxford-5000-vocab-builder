package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/AkifhanIlgaz/vocab-builder/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Auth *auth.Client
}

type UserService struct {
	Collection *mongo.Collection
}

func NewUserService(client *mongo.Client) UserService {
	return UserService{
		Collection: client.Database(Database).Collection(UsersCollection),
	}
}

type User struct {
	Uid          primitive.ObjectID `json:"uid" bson:"_id"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	CreatedAt    time.Time          `json:"-" bson:"createdAt"`
}

func (service *UserService) Create(email, password string) (*User, error) {
	// Check if user exists
	isUserExist, err := service.CheckIfUserExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if isUserExist {
		return nil, fmt.Errorf("create user: %w", errors.ErrEmailTaken)
	}

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

	_, err = service.Collection.InsertOne(context.TODO(), user)

	if err != nil {
		return nil, errors.MongoError(fmt.Errorf("create user: %w", err))
	}

	return &user, nil
}

func (service *UserService) CheckIfUserExistsByEmail(email string) (bool, error) {
	count, err := service.Collection.CountDocuments(context.TODO(), bson.M{
		"email": email,
	})
	if err != nil {
		return false, errors.MongoError(fmt.Errorf("check user by email: %w", err))
	}

	return count > 0, nil
}
