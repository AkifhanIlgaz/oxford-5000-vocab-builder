package models

import (
	"context"
	"fmt"
	"time"

	"github.com/AkifhanIlgaz/vocab-builder/errors"
	"github.com/AkifhanIlgaz/vocab-builder/rand"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultAccessTokenExpireDuration = 60 * time.Minute

type TokenService struct {
	UsersCollection           *mongo.Collection
	RefreshTokenCollection    *mongo.Collection
	accessTokenExpireDuration time.Duration
}

func NewTokenService(client *mongo.Client, accessTokenExpireDuration *time.Duration) TokenService {
	expireDuration := defaultAccessTokenExpireDuration
	if accessTokenExpireDuration != nil {
		expireDuration = *accessTokenExpireDuration
	}

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
		accessTokenExpireDuration: expireDuration,
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

type refreshTokenInfo struct {
	Uid          string `bson:"uid"`
	RefreshToken string `bson:"refreshToken"`
}

func (service *TokenService) NewRefreshToken(uid string) (string, error) {
	// ! We don't need to check whether user exists since we are checking it when client wants to access user information

	refreshToken, err := rand.String(32)
	if err != nil {
		return "", fmt.Errorf("new refresh token: %w", err)
	}

	_, err = service.RefreshTokenCollection.InsertOne(context.TODO(), refreshTokenInfo{
		Uid:          uid,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", fmt.Errorf("new access token: %w", err)
	}

	return refreshToken, nil
}

func (service *TokenService) DeleteRefreshToken(uid string) error {
	res, err := service.RefreshTokenCollection.DeleteOne(context.TODO(), bson.M{
		"uid": uid,
	})
	if err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	if res.DeletedCount == 0 {
		return errors.New("refresh token doesn't exist")
	}

	return nil
}

func (service *TokenService) RefreshAccessToken(refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	var info refreshTokenInfo

	newRefreshToken, err = rand.String(32)
	if err != nil {
		return "", "", fmt.Errorf("refresh access token: %w", err)
	}

	err = service.RefreshTokenCollection.FindOneAndUpdate(context.TODO(), bson.M{
		"refreshToken": refreshToken,
	}, bson.M{
		"refreshToken": newRefreshToken,
	}).Decode(&info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", errors.New("refresh access token: refresh token doesn't exist for the user:")
		}
		return "", "", fmt.Errorf("refresh access token: %w", err)
	}

	newAccessToken, err = service.NewAccessToken(info.Uid)
	if err != nil {
		return "", "", fmt.Errorf("refresh access token: %w", err)
	}

	return
}

func (service *TokenService) ParseAccessToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
}

func keyFunc(t *jwt.Token) (interface{}, error) {
	return Secret, nil
}
