package router

import (
	"context"
	"fmt"
	"sort"

	"github.com/mbadhala/go-common-api-framework/core"
	"github.com/mbadhala/go-common-api-framework/middleware"
)

type route struct {
	method      string
	path        string
	handler     core.Handler
	middlewares []middleware.Middleware
	meta        *RouteMeta
}

type Router struct {
	parent      *Router
	routes      []route
	middlewares []middleware.Middleware
	prefix      string
}

func New() *Router {
	return &Router{}
}

func (r *Router) Use(m ...middleware.Middleware) {
	r.middlewares = append(r.middlewares, m...)
}

func (r *Router) Group(prefix string) *Router {
	return &Router{
		parent: r,
		prefix: r.prefix + prefix,
	}
}

func (r *Router) root() *Router {
	if r.parent == nil {
		return r
	}
	return r.parent.root()
}

func (r *Router) collectMiddleware() []middleware.Middleware {
	if r.parent == nil {
		return r.middlewares
	}
	return append(r.parent.collectMiddleware(), r.middlewares...)
}

func (r *Router) Handle(method, path string, h core.Handler, m ...middleware.Middleware) {
	allMiddleware := append(r.collectMiddleware(), m...)

	rt := route{
		method:      method,
		path:        r.prefix + path,
		handler:     h,
		middlewares: allMiddleware,
	}

	root := r.root()
	root.routes = append(root.routes, rt)
}

func (r *Router) HandleWithMeta(
	method, path string, h core.Handler, meta RouteMeta, m ...middleware.Middleware) {
	meta.Method = method
	meta.Path = r.prefix + path
	allMiddleware := append(r.collectMiddleware(), m...)

	rt := route{
		method:      method,
		path:        r.prefix + path,
		handler:     h,
		middlewares: allMiddleware,
		meta:        &meta,
	}

	root := r.root()
	root.routes = append(root.routes, rt)
}

func (r *Router) Serve(ctx context.Context, req core.Request) core.Response {
	for _, rt := range r.root().routes {
		if rt.method == req.Method && rt.path == req.Path {
			h := middleware.Chain(rt.middlewares...)(rt.handler)
			return h(ctx, req)
		}
	}
	return core.Error(404, "not found")
}

func (r *Router) Routes() []struct {
	Method string
	Path   string
} {
	root := r.root()

	var out []struct {
		Method string
		Path   string
	}

	for _, rt := range root.routes {
		out = append(out, struct {
			Method string
			Path   string
		}{
			Method: rt.method,
			Path:   rt.path,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Path == out[j].Path {
			return out[i].Method < out[j].Method
		}
		return out[i].Path < out[j].Path
	})

	return out
}

func (r *Router) PrintRoutes() {
	fmt.Println("==== Registered Routes ====")
	for _, rt := range r.Routes() {
		fmt.Printf("%-6s %s\n", rt.Method, rt.Path)
	}
	fmt.Println("===========================")
}
