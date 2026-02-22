package adapter

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mbadhala/go-common-api-framework/core"
	"github.com/mbadhala/go-common-api-framework/router"
)

func Gin(r *router.Router) *gin.Engine {
	g := gin.New()

	g.Any("/*path", func(c *gin.Context) {

		req := core.Request{
			Method:      c.Request.Method,
			Path:        c.Request.URL.Path,
			Headers:     map[string]string{},
			QueryParams: map[string]string{},
			PathParams:  map[string]string{},
			Form:        map[string]string{},
			Files:       map[string]core.File{},
		}

		// Headers
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				req.Headers[k] = v[0]
			}
		}

		// Query params
		for k, v := range c.Request.URL.Query() {
			if len(v) > 0 {
				req.QueryParams[k] = v[0]
			}
		}

		// Path params
		for _, p := range c.Params {
			req.PathParams[p.Key] = p.Value
		}

		// Detect multipart
		if c.ContentType() == "multipart/form-data" {

			err := c.Request.ParseMultipartForm(32 << 20)
			if err == nil {

				// Form fields
				for k, v := range c.Request.MultipartForm.Value {
					if len(v) > 0 {
						req.Form[k] = v[0]
					}
				}

				// Files
				for k, files := range c.Request.MultipartForm.File {
					if len(files) > 0 {
						fh := files[0]

						f, err := fh.Open()
						if err == nil {
							content, _ := io.ReadAll(f)
							f.Close()

							req.Files[k] = core.File{
								Filename: fh.Filename,
								Size:     fh.Size,
								Content:  content,
							}
						}
					}
				}
			}

		} else {
			// Normal body
			body, _ := io.ReadAll(c.Request.Body)
			req.Body = body
		}

		// Call core router
		res := r.Serve(c.Request.Context(), req)

		// Write response
		for k, v := range res.Headers {
			c.Writer.Header().Set(k, v)
		}

		c.Data(res.StatusCode, http.DetectContentType(res.Body), res.Body)
	})

	return g
}
