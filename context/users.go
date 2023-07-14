package context

import (
	"context"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type key string

const userKey key = "user"

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)

	if user, ok := val.(*models.User); ok {
		return user
	}

	return nil
}
