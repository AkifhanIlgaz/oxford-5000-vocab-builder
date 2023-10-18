package models

import (
	"context"
	"fmt"
	"time"

	"github.com/AkifhanIlgaz/vocab-builder/errors"
	"github.com/AkifhanIlgaz/vocab-builder/rand"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	type data struct {
		Uid          string `bson:"uid"`
		RefreshToken string `bson:"refreshToken"`
	}

	_, err = service.RefreshTokenCollection.InsertOne(context.TODO(), data{
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

	if isMatch := checkRefreshClaims(&refreshClaims, uid, refreshToken); !isMatch {
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

func (service *TokenService) ParseAccessToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
}

func keyFunc(t *jwt.Token) (interface{}, error) {
	return Secret, nil
}
