package core

type Request struct {
	Method      string
	Path        string
	Headers     map[string]string
	Body        []byte
	PathParams  map[string]string
	QueryParams map[string][]string
}
