package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenService struct {
	Collection *mongo.Collection
}

func NewTokenService(client *mongo.Client) TokenService {
	collection := client.Database(Database).Collection(RefreshTokenCollection)

	
}

type UserClaims struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}

func NewIdToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpireDuration)),
		},
	})

	t, err := token.SignedString(Secret)
	if err != nil {
		return "", fmt.Errorf("new id token: %w", err)
	}

	return t, nil
}
