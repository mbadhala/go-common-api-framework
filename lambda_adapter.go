package revinexgorestapiframework

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func (a *App) LambdaHandler() func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		for _, route := range a.routes {
			if route.Path == req.Path && route.Method == req.HTTPMethod {

				r := Request{
					Method:  req.HTTPMethod,
					Path:    req.Path,
					Headers: req.Headers,
					Body:    []byte(req.Body),
				}

				resp, err := route.Handler(ctx, r)
				if err != nil {
					return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
				}

				b, _ := json.Marshal(resp.Body)

				return events.APIGatewayProxyResponse{
					StatusCode: resp.StatusCode,
					Headers:    resp.Headers,
					Body:       string(b),
				}, nil
			}
		}

		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "not found"}, nil
	}
}
