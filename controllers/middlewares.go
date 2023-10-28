package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/AkifhanIlgaz/vocab-builder/errors"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/AkifhanIlgaz/vocab-builder/oauth"
)

const BearerScheme = "Bearer "

type AccessTokenMiddleware struct {
	TokenService *models.TokenService
	GoogleOauth  *oauth.GoogleOAuth
}

func NewAccessTokenMiddleware(tokenService *models.TokenService, googleOauth *oauth.GoogleOAuth) *AccessTokenMiddleware {
	return &AccessTokenMiddleware{
		TokenService: tokenService,
		GoogleOauth:  googleOauth,
	}
}

func (middleware *AccessTokenMiddleware) AccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provider := r.Header.Get("provider")

		if provider == "email" {
			middleware.withEmailProvider(next, w, r)
		} else if provider == "google" {
			middleware.withGoogleProvider(next, w, r)
		}

	})
}

func (middleware *AccessTokenMiddleware) withEmailProvider(next http.Handler, w http.ResponseWriter, r *http.Request) {
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
		next.ServeHTTP(w, r)
		return
	}

	ctx := r.Context()
	ctx = context.WithUid(ctx, uid)
	r = r.WithContext(ctx)

	next.ServeHTTP(w, r)

}

func (middleware *AccessTokenMiddleware) withGoogleProvider(next http.Handler, w http.ResponseWriter, r *http.Request) {
	// Extract access token from Authorization Header
	tokenString, err := parseBearer(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Printf("access token middleware: %v\n", err)
		next.ServeHTTP(w, r)
		return
	}

	req, err := http.NewRequest("GET", middleware.GoogleOauth.UserUrl, nil)
	if err != nil {
		fmt.Println(err)
		next.ServeHTTP(w, r)
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// TODO: Return with token is expired error
		fmt.Println(err)
		return
	}

	body, _ := io.ReadAll(res.Body)
	var respBody map[string]string
	json.Unmarshal(body, &respBody)

	uid, ok := respBody["sub"]
	if !ok {
		fmt.Println("sub claim doesn't exist")
		next.ServeHTTP(w, r)
		return
	}

	ctx := r.Context()
	ctx = context.WithUid(ctx, uid)
	r = r.WithContext(ctx)

	next.ServeHTTP(w, r)

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
