package models

import (
	"context"
	"fmt"
	"time"

	"github.com/AkifhanIlgaz/vocab-builder/rand"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenService struct {
	UsersCollection        *mongo.Collection
	RefreshTokenCollection *mongo.Collection
	idTokenExpireDuration  time.Duration
}

func NewTokenService(client *mongo.Client, idTokenExpireDuration time.Duration) TokenService {
	return TokenService{
		UsersCollection:        getCollection(client, UsersCollection),
		RefreshTokenCollection: getCollection(client, RefreshTokenCollection),
		idTokenExpireDuration:  idTokenExpireDuration,
	}
}

type UserClaims struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}

func (service *TokenService) NewIdToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.idTokenExpireDuration)),
		},
	})

	t, err := token.SignedString(Secret)
	if err != nil {
		return "", fmt.Errorf("new id token: %w", err)
	}

	return t, nil
}

type RefreshClaims struct {
	Uid          string
	RefreshToken string
	jwt.RegisteredClaims
}

// TODO: Create Refresh Token for user
/*
	!Errors
		User doesn't exist
*/
func (service *TokenService) NewRefreshToken(uid string) (string, error) {
	// TODO: Check if user exists
	exists, err := service.CheckIfUserExists(uid)
	if err != nil {
		return "", fmt.Errorf("new refresh token: %w", err)
	}

	if exists {
		refreshToken, err := rand.String(32)
		if err != nil {
			return "", fmt.Errorf("new refresh token: %w", err)
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims{
			Uid:          uid,
			RefreshToken: refreshToken,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(time.Now()),
			},
		})

		t, err := token.SignedString(Secret)
		if err != nil {
			return "", fmt.Errorf("new id token: %w", err)
		}

		return t, nil
	} else {
		return "", fmt.Errorf("new refresh token: user doesn't exist")
	}

}

/*
!Errors
User doesn't exist
Refresh token expired && not valid => refresh token isn't same as the refresh token on DB
*/
func (service *TokenService) RefreshIdToken(uid, refreshToken string) (string, error) {

	panic("Implement")
}

// Returns true if user exists
func (service *TokenService) CheckIfUserExists(uid string) (bool, error) {
	objId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return false, fmt.Errorf("check if user exists: %w", err)
	}

	count, err := service.UsersCollection.CountDocuments(context.TODO(), bson.M{
		"_id": objId,
	})

	if err != nil {
		return false, fmt.Errorf("check if user exists: %w", err)
	}

	fmt.Println("count", count)

	return count > 0, nil
}
