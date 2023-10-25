package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/vocab-builder/context"
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
		// Extract access token from Authorization Header
		tokenString, err := parseBearer(r.Header.Get("Authorization"))
		if err != nil {
			fmt.Printf("access token middleware: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		// Parse access token
		token, err := middleware.TokenService.ParseAccessToken(tokenString)
		if err != nil {
			fmt.Println(w, "parse access token", err)
			next.ServeHTTP(w, r)
			return
		}

		uid, err := token.Claims.GetSubject()
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println(uid)

		ctx := r.Context()
		ctx = context.WithUid(ctx, uid)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
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

	return splitAuth, nil
}
