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

func NewTokenService(client *mongo.Client) *TokenService {
	collection := client.Database(Database).Collection(RefreshTokenCollection)

	return &TokenService{
		Collection: collection,
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpireDuration)),
		},
	})

	t, err := token.SignedString(Secret)
	if err != nil {
		return "", fmt.Errorf("new id token: %w", err)
	}

	return t, nil
}

// TODO: Create Refresh Token for user
/*
	!Errors
		User doesn't exist

*/
func (service *TokenService) NewRefreshToken(uid string) (string, error) {
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

/*
!Errors
User doesn't exist
Refresh token expired && not valid => refresh token isn't same as the refresh token on DB
*/
func (service *TokenService) RefreshIdToken(uid, refreshToken string) (string, error) {

	panic("Implement")
}
