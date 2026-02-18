package revinexgorestapiframework

import "context"

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       any
}

type HandlerFunc func(ctx context.Context, req Request) (Response, error)
type Middleware func(HandlerFunc) HandlerFunc

type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
}
