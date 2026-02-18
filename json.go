package revinexgorestapiframework

import (
	"context"
	"encoding/json"
)

func JSON(status int, v any) Response {
	return Response{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: v,
	}
}

type TypedHandler[I any, O any] func(ctx context.Context, input I) (O, error)

func Wrap[I any, O any](h TypedHandler[I, O]) HandlerFunc {
	return func(ctx context.Context, req Request) (Response, error) {
		var input I

		if err := json.Unmarshal(req.Body, &input); err != nil {
			return JSON(400, map[string]string{"error": "invalid json"}), nil
		}

		output, err := h(ctx, input)
		if err != nil {
			return JSON(500, map[string]string{"error": err.Error()}), nil
		}

		return JSON(200, output), nil
	}
}
