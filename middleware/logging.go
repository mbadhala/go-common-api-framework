package middleware

import (
	"context"
	"log"
	"time"

	"github.com/mbadhala/go-common-api-framework/core"
)

func Logging() Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx context.Context, req core.Request) core.Response {
			start := time.Now()
			res := next(ctx, req)
			log.Printf("%s %s %d %v",
				req.Method,
				req.Path,
				res.StatusCode,
				time.Since(start),
			)
			return res
		}
	}
}
