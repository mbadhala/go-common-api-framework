package revinexgorestapiframework

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) MountGin(g *gin.Engine) {
	for _, route := range a.routes {
		switch route.Method {
		case "GET":
			g.GET(route.Path, buildGinHandler(route.Handler))
		case "POST":
			g.POST(route.Path, buildGinHandler(route.Handler))
		}
	}
}

func buildGinHandler(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		body, _ := io.ReadAll(c.Request.Body)

		headers := map[string]string{}
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}

		req := Request{
			Method:  c.Request.Method,
			Path:    c.FullPath(),
			Headers: headers,
			Body:    body,
		}

		resp, err := h(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for k, v := range resp.Headers {
			c.Header(k, v)
		}

		c.JSON(resp.StatusCode, resp.Body)
	}
}
