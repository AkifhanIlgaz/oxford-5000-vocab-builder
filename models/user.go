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
