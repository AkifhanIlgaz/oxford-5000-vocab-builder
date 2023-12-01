package errors

import (
	"errors"
	"net/http"
)

var (
	ErrEmailTaken    error = errors.New("email is taken")
	ErrUserNotExist  error = errors.New("user does not exist")
	ErrWrongPassword error = errors.New("wrong password")
	ErrMongo         error = errors.New("something went wrong with MongoDB")
)

var (
	Join = errors.Join
	As   = errors.As
	Is   = errors.Is
	New  = errors.New
)

func MongoError(err error) error {
	return Join(ErrMongo, err)
}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}
