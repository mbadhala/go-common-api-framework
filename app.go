package revinexgorestapiframework

type App struct {
	routes      []Route
	middlewares []Middleware
}

func New() *App {
	return &App{}
}

func (a *App) Use(m Middleware) {
	a.middlewares = append(a.middlewares, m)
}

func (a *App) Handle(method, path string, h HandlerFunc) {
	for i := len(a.middlewares) - 1; i >= 0; i-- {
		h = a.middlewares[i](h)
	}

	a.routes = append(a.routes, Route{
		Method:  method,
		Path:    path,
		Handler: h,
	})
}

func (a *App) Routes() []Route {
	return a.routes
}
