package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/platform/framework/core"
)

func RequestTiming() Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx context.Context, req core.Request) core.Response {
			start := time.Now()

			if req.Headers == nil {
				req.Headers = map[string]string{}
			}
			req.Headers["X-Request-Received"] = start.Format(time.RFC3339Nano)

			// Call next handler
			res := next(ctx, req)

			if res.Headers == nil {
				res.Headers = map[string]string{}
			}
			res.Headers["X-Request-Completed"] = time.Now().Format(time.RFC3339Nano)
			res.Headers["X-Request-Received"] = start.Format(time.RFC3339Nano)

			latency := time.Since(start)
			res.Headers["X-Execution-latency"] = latency.String()
			fmt.Printf("[RequestTiming] %s %s took %v\n", req.Method, req.Path, latency)

			return res
		}
	}
}
