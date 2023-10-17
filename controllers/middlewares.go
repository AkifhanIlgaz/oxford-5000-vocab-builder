package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/vocab-builder/errors"
	"github.com/AkifhanIlgaz/vocab-builder/models"
)

const BearerScheme = "Bearer "

type AccessTokenMiddleware struct {
	TokenService *models.TokenService
}

func NewAccessTokenMiddleware(tokenService *models.TokenService) *AccessTokenMiddleware {
	return &AccessTokenMiddleware{
		TokenService: tokenService,
	}
}

func (middleware *AccessTokenMiddleware) AccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := parseBearer(r.Header.Get("Authorization"))
		if err != nil {
			fmt.Printf("access token middleware: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		token, err := middleware.TokenService.ParseAccessToken(tokenString)
		if err != nil {
			fmt.Fprint(w, "parse id token")
			return
		}

		enc := json.NewEncoder(w)
		err = enc.Encode(token.Claims)
		if err != nil {
			fmt.Println(err)
		}
	})
}

func parseBearer(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	splitAuth, found := strings.CutPrefix(authHeader, BearerScheme)
	if !found {
		return "", errors.New("invalid bearer scheme")
	}

	// 32 byte => base64 encoding => 44
	if len(splitAuth) != 44 {
		return "", errors.New("invalid bearer scheme")
	}

	return splitAuth, nil
}
