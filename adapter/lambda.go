package adapter

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/yourorg/platform/framework/core"
	"github.com/yourorg/platform/framework/router"
)

func Lambda(r *router.Router) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		req := core.Request{
			Method:  e.HTTPMethod,
			Path:    e.Path,
			Headers: e.Headers,
			Body:    []byte(e.Body),
		}

		res := r.Serve(ctx, req)

		return events.APIGatewayProxyResponse{
			StatusCode: res.StatusCode,
			Headers:    res.Headers,
			Body:       string(res.Body),
		}, nil
	}
}
