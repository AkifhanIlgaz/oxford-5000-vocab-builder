package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"github.com/AkifhanIlgaz/vocab-builder/rand"
)

const bytesPerToken = 32

type Session struct {
	Id        int
	UserId    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *firebase.App
}

func (service *SessionService) Create(userId int) (*Session, error) {
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	session := Session{
		UserId:    userId,
		Token:     token,
		TokenHash: service.hash(token),
	}

	row := service.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash) 
		VALUES ( $1, $2 ) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id;
	`, session.UserId, session.TokenHash)

	err = row.Scan(&session.Id)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &session, nil
}

func (service *SessionService) User(token string) (*User, error) {
	tokenHash := service.hash(token)
	var user User

	row := service.DB.QueryRow(`
		SELECT users.id,
				users.email,
				users.password_hash
		FROM sessions 
		JOIN users ON users.id = sessions.user_id
		WHERE sessions.token_hash = $1;
	`, tokenHash)

	err := row.Scan(&user.Id, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

func (service *SessionService) Delete(token string) error {
	tokenHash := service.hash(token)

	_, err := service.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1;
	`, tokenHash)

	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (service *SessionService) hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}
