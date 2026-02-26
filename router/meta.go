package router

type RouteMeta struct {
	Method       string
	Path         string
	Summary      string
	Description  string
	RequestType  any
	ResponseType any
}
