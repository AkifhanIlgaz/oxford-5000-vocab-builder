package controllers

import (
	"strings"

	"github.com/AkifhanIlgaz/vocab-builder/errors"
)

const BearerScheme = "Bearer "

func ParseBearer(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	splitAuth, found := strings.CutPrefix(authHeader, BearerScheme)
	if !found {
		return "", errors.New("invalid bearer scheme")
	}

	if len(splitAuth) != 44 {
		return "", errors.New("invalid bearer scheme")
	}

	return splitAuth, nil
}
