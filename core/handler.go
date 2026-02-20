package core

import (
	"context"
	"encoding/json"
)

type Handler func(ctx context.Context, req Request) Response

func Wrap[I any, O any](fn func(context.Context, Request, I) (O, error)) Handler {
	return func(ctx context.Context, req Request) Response {
		var input I
		if err := json.Unmarshal(req.Body, &input); err != nil {
			return Error(400, "invalid json")
		}

		out, err := fn(ctx, req, input)
		if err != nil {
			return Error(500, err.Error())
		}

		return JSON(200, out)
	}
}

func WrapNoBody[O any](fn func(context.Context, Request) (O, error)) Handler {
	return func(ctx context.Context, req Request) Response {
		out, err := fn(ctx, req)
		if err != nil {
			return Error(500, err.Error())
		}
		return JSON(200, out)
	}
}
