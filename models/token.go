package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenService struct {
	UsersCollection            *mongo.Collection
	RefreshTokenCollection     *mongo.Collection
	idTokenExpireDuration      time.Duration
	refreshTokenExpireDuration time.Duration
}

func NewTokenService(client *mongo.Client, idTokenExpireDuration, refreshTokenExpireDuration time.Duration) *TokenService {
	return &TokenService{
		UsersCollection:            getCollection(client, UsersCollection),
		RefreshTokenCollection:     getCollection(client, RefreshTokenCollection),
		idTokenExpireDuration:      idTokenExpireDuration,
		refreshTokenExpireDuration: refreshTokenExpireDuration,
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.refreshTokenExpireDuration)),
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
	// TODO: Create random string as

	panic("Implement")
}
