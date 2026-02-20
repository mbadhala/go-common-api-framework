package middleware

import (
	"github.com/mbadhala/go-common-api-framework/core"
)

type Middleware func(core.Handler) core.Handler

func Chain(m ...Middleware) Middleware {
	return func(next core.Handler) core.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
