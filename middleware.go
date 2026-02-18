package revinexgorestapiframework

import (
	"context"
	"log"
	"time"
)

func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (Response, error) {
		start := time.Now()
		resp, err := next(ctx, req)
		log.Printf("%s %s %d %v", req.Method, req.Path, resp.StatusCode, time.Since(start))
		return resp, err
	}
}

func AuthMiddleware(secret string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, req Request) (Response, error) {

			token := req.Headers["Authorization"]
			if token != secret {
				return JSON(401, map[string]string{"error": "unauthorized"}), nil
			}

			return next(ctx, req)
		}
	}
}
