package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mbadhala/go-common-api-framework/core"
	"github.com/mbadhala/go-common-api-framework/middleware"
)

type ctxKey struct{}

var claimsKey ctxKey

type Claims struct {
	Sub   string   `json:"sub"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func JWT(publicKey *rsa.PublicKey) middleware.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx context.Context, req core.Request) core.Response {

			authHeader := req.Headers["Authorization"]
			if authHeader == "" {
				return core.Error(401, "missing token")
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return core.Error(401, "invalid authorization header")
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
				return publicKey, nil
			})

			if err != nil || !token.Valid {
				return core.Error(401, "invalid token")
			}

			claims := token.Claims.(*Claims)

			ctx = context.WithValue(ctx, claimsKey, claims)

			return next(ctx, req)
		}
	}
}

func GetClaims(ctx context.Context) (*Claims, error) {
	v := ctx.Value(claimsKey)
	if v == nil {
		return nil, errors.New("no claims in context")
	}

	claims, ok := v.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	return claims, nil
}

func Policy(fn func(*Claims) bool) middleware.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx context.Context, req core.Request) core.Response {

			claims, err := GetClaims(ctx)
			if err != nil {
				return core.Error(401, "unauthorized")
			}

			if !fn(claims) {
				return core.Error(403, "forbidden")
			}

			return next(ctx, req)
		}
	}
}
