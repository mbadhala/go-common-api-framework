package core

import "encoding/json"

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func JSON(status int, v any) Response {
	b, _ := json.Marshal(v)
	return Response{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: b,
	}
}

func Error(status int, msg string) Response {
	return JSON(status, map[string]string{
		"error": msg,
	})
}
