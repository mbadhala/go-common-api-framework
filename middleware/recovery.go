package middleware

import (
	"context"
	"log"

	"github.com/yourorg/platform/framework/core"
)

func Recovery() Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx context.Context, req core.Request) (res core.Response) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("panic:", r)
					res = core.Error(500, "internal error")
				}
			}()
			return next(ctx, req)
		}
	}
}
