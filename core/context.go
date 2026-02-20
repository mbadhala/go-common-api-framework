package core

import "context"

type ctxKey string

const claimsKey ctxKey = "claims"

func WithValue(ctx context.Context, key string, val any) context.Context {
	return context.WithValue(ctx, ctxKey(key), val)
}

func GetValue[T any](ctx context.Context, key string) (T, bool) {
	v := ctx.Value(ctxKey(key))
	if v == nil {
		var zero T
		return zero, false
	}
	val, ok := v.(T)
	return val, ok
}

func setClaims(ctx context.Context, claims any) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func getClaims(ctx context.Context) any {
	return ctx.Value(claimsKey)
}
