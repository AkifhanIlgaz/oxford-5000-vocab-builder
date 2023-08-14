package context

import (
	"context"
)

type key string

const uidKey key = "uid"

func WithUid(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, uidKey, uid)
}

func Uid(ctx context.Context) string {
	if uid, ok := ctx.Value(uidKey).(string); ok {
		return uid
	}
	return ""
}
