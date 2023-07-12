package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (service *UserService) Create(email, password string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	row := service.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES (
			$1,
			$2
		) 
		RETURNING id;
	`, user.Email, user.PasswordHash)

	err = row.Scan(&user.Id)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, fmt.Errorf("email is taken")
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (service *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	user := User{
		Email: email,
	}

	row := service.DB.QueryRow(`
		SELECT id, password_hash FROM users
		WHERE email = $1;
	`, user.Email)

	err := row.Scan(&user.Id, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate user: %w", err)
	}

	return &user, nil
}

func (service *UserService) UpdatePassword(userId int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)

	_, err = service.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1;
	`, userId, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}
