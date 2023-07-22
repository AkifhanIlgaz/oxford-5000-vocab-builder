package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/AkifhanIlgaz/vocab-builder/rand"
)

const ResetDuration = 1 * time.Hour

type PasswordReset struct {
	Id        int
	UserId    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
}

// Create a password reset token for the given email
func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	// Check if the user exists by using e mail
	email = strings.TrimSpace(strings.ToLower(email))

	var userId int
	row := service.DB.QueryRow(`
		SELECT FROM users
		WHERE email = $1
		RETURNING id;
	`, email)
	if err := row.Scan(&userId); err != nil {
		return nil, fmt.Errorf("create password reset token: %w", err)
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create password reset token: %w", err)
	}

	passwordReset := PasswordReset{
		UserId:    userId,
		Token:     token,
		TokenHash: service.hash(token),
		ExpiresAt: time.Now().Add(ResetDuration),
	}

	row = service.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;
	`, userId, passwordReset.TokenHash, passwordReset.ExpiresAt)

	if err := row.Scan(&passwordReset.Id); err != nil {
		return nil, fmt.Errorf("create password reset token: %w", err)
	}

	return &passwordReset, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := service.hash(token)
	var user User
	var passwordReset PasswordReset

	row := service.DB.QueryRow(`
		SELECT  password_resets.id,
				password_resets.expires_at,
				users.id,
				users.email,
				users.password_hash
		FROM password_resets
				JOIN users on users.id = password_reset.user_id
		WHERE password_resets.token_hash = $1;
	`, tokenHash)

	err := row.Scan(&passwordReset.Id, &passwordReset.ExpiresAt, &user.Id, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("consume reset token: %w", err)
	}

	if time.Now().After(passwordReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}

	err = service.delete(passwordReset.Id)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (service *PasswordResetService) hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (service *PasswordResetService) delete(id int) error {
	if _, err := service.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1;
	`, id); err != nil {
		return fmt.Errorf("delete reset token: %w", err)
	}

	return nil
}
