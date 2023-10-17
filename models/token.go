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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenService struct {
	UsersCollection           *mongo.Collection
	RefreshTokenCollection    *mongo.Collection
	accessTokenExpireDuration time.Duration
}

func NewTokenService(client *mongo.Client, accessTokenExpireDuration time.Duration) TokenService {
	usersCollection := getCollection(client, UsersCollection)
	refreshTokenCollection := getCollection(client, RefreshTokenCollection)

	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"uid": 2},
		Options: options.Index().SetUnique(true),
	}

	// TODO: Check error
	refreshTokenCollection.Indexes().CreateOne(context.TODO(), indexModel)

	return TokenService{
		UsersCollection:           usersCollection,
		RefreshTokenCollection:    refreshTokenCollection,
		accessTokenExpireDuration: accessTokenExpireDuration,
	}
}

func (service *TokenService) NewAccessToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject:   uid,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.accessTokenExpireDuration)),
		},
	)

	t, err := token.SignedString(Secret)
	if err != nil {
		return "", fmt.Errorf("new access token: %w", err)
	}

	return t, nil
}

type RefreshClaims struct {
	Uid                  string `bson:"uid"`
	RefreshToken         string `bson:"refreshToken"`
	jwt.RegisteredClaims `bson:"-"`
}

func (service *TokenService) NewRefreshToken(uid string) (string, error) {
	exists, err := service.CheckIfUserExists(uid)
	if err != nil {
		return "", fmt.Errorf("new refresh token: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("new refresh token: user doesn't exist")

	}

	refreshToken, err := rand.String(32)
	if err != nil {
		return "", fmt.Errorf("new refresh token: %w", err)
	}

	claims := RefreshClaims{
		Uid:          uid,
		RefreshToken: refreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	_, err = service.RefreshTokenCollection.InsertOne(context.TODO(), claims)
	if err != nil {
		return "", fmt.Errorf("new access token: %w", err)
	}

	t, err := token.SignedString(Secret)
	if err != nil {
		return "", fmt.Errorf("new access token: %w", err)
	}

	return t, nil

}

/*
!Errors
User doesn't exist
Refresh token expired && not valid => refresh token isn't same as the refresh token on DB
*/
func (service *TokenService) RefreshAccessToken(uid, refreshToken string) (string, error) {
	exists, err := service.CheckIfUserExists(uid)
	if err != nil {
		return "", fmt.Errorf("refresh id token: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("refresh id token: user doesn't exist")
	}

	var refreshClaims RefreshClaims

	err = service.RefreshTokenCollection.FindOne(context.TODO(), bson.M{
		"uid": uid,
	}).Decode(&refreshClaims)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("refresh id token: refresh token doesn't exist for the user: %v", uid)
		}
		return "", fmt.Errorf("refresh id token: %w", err)
	}

	isMatch := checkRefreshClaims(&refreshClaims, uid, refreshToken)
	if !isMatch {
		return "", fmt.Errorf("invalid refresh token for the user")
	}

	token, err := service.NewAccessToken(uid)
	if err != nil {
		return "", fmt.Errorf("refresh id token: %w", err)
	}

	return token, nil
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

	return count > 0, nil
}

func checkRefreshClaims(claims *RefreshClaims, uid, refreshToken string) bool {
	return claims.Subject == uid && claims.RefreshToken == refreshToken
}

func (service *TokenService) ParseAccessToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
}

func (service *TokenService) ParseRefreshToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
}

func keyFunc(t *jwt.Token) (interface{}, error) {
	return Secret, nil
}
