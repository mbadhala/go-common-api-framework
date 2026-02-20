package adapter

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/platform/framework/core"
	"github.com/yourorg/platform/framework/router"
)

func Gin(r *router.Router) *gin.Engine {
	g := gin.New()

	g.Any("/*path", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)

		req := core.Request{
			Method:  c.Request.Method,
			Path:    c.Request.URL.Path,
			Headers: map[string]string{},
			Body:    body,
		}

		for k, v := range c.Request.Header {
			req.Headers[k] = v[0]
		}

		res := r.Serve(c.Request.Context(), req)

		for k, v := range res.Headers {
			c.Writer.Header().Set(k, v)
		}

		c.Data(res.StatusCode, "application/json", res.Body)
	})

	return g
}
