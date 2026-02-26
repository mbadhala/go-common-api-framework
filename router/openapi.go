package router

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type OpenAPI struct {
	OpenAPI string                 `json:"openapi"`
	Info    map[string]string      `json:"info"`
	Paths   map[string]interface{} `json:"paths"`
}

func GenerateOpenAPI(r *Router) ([]byte, error) {

	doc := OpenAPI{
		OpenAPI: "3.0.3",
		Info: map[string]string{
			"title":   "API Service",
			"version": "1.0.0",
		},
		Paths: map[string]interface{}{},
	}

	for _, route := range r.routes {

		if route.meta == nil {
			continue
		}

		method := strings.ToLower(route.method)

		if doc.Paths[route.path] == nil {
			doc.Paths[route.path] = map[string]interface{}{}
		}

		pathItem := doc.Paths[route.path].(map[string]interface{})

		// -------- RESPONSE FIX ----------
		responseObj := map[string]interface{}{
			"description": "Success",
		}

		if route.meta.ResponseType != nil {
			responseObj["content"] = map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": generateSchema(route.meta.ResponseType),
				},
			}
		}

		// -------- OPERATION ----------
		operation := map[string]interface{}{
			"summary":     route.meta.Summary,
			"description": route.meta.Description,
			"responses": map[string]interface{}{
				"200": responseObj,
			},
		}

		// -------- REQUEST ----------
		if route.meta.RequestType != nil {
			operation["requestBody"] = map[string]interface{}{
				"required": true,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": generateSchema(route.meta.RequestType),
					},
				},
			}
		}
		fmt.Print(operation)
		pathItem[method] = operation
	}

	return json.MarshalIndent(doc, "", "  ")
}

func generateSchema(model any) map[string]interface{} {

	t := reflect.TypeOf(model)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	properties := map[string]interface{}{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		properties[jsonTag] = map[string]string{
			"type": mapType(field.Type.Kind()),
		}
	}

	return map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
}

func mapType(k reflect.Kind) string {
	switch k {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int64:
		return "integer"
	case reflect.Bool:
		return "boolean"
	default:
		return "string"
	}
}
