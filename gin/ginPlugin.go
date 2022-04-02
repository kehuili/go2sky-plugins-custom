package ginPlugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/gin-gonic/gin"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const componentIDGINHttpServer = 5006

//Middleware gin middleware return HandlerFunc  with tracing.
func Middleware(engine *gin.Engine, tracer *go2sky.Tracer, opts ...Option) gin.HandlerFunc {
	if engine == nil || tracer == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	options := &options{}
	for _, o := range opts {
		o(options)
	}

	return func(c *gin.Context) {
		// ignore designated paths
		for _, path := range options.excludePaths {
			if c.FullPath() == path {
				c.Next()
				return
			}
		}

		span, ctx, err := tracer.CreateEntrySpan(c.Request.Context(), getOperationName(c), func(key string) (string, error) {
			sw := c.Request.Header.Get(key)
			if sw != "" {
				return sw, nil
			}
			// no header, try body
			// rs := make(map[string]interface{}, 0)
			// if err := c.ShouldBindBodyWith(&rs, binding.JSON); err != nil {
			// 	return "", nil
			// }

			for _, path := range options.fromBodyPaths {
				if c.FullPath() == path {
					body, _ := ioutil.ReadAll(c.Request.Body)
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
					rs := make(map[string]interface{})
					if err := json.Unmarshal(body, &rs); err != nil {
						return "", nil
					}
					if rs["swHeaders"] != nil && rs["swHeaders"].(map[string]interface{})[key] != nil {
						sw = rs["swHeaders"].(map[string]interface{})[key].(string)
					}
					return sw, nil
				}
			}
			return sw, nil
		})
		if err != nil {
			c.Next()
			return
		}
		span.SetComponent(componentIDGINHttpServer)
		span.Tag(go2sky.TagHTTPMethod, c.Request.Method)
		span.Tag(go2sky.TagURL, c.Request.Host+c.Request.URL.Path)
		span.SetSpanLayer(agentv3.SpanLayer_Http)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if len(c.Errors) > 0 {
			span.Error(time.Now(), c.Errors.String())
		}
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(c.Writer.Status()))
		span.End()
	}
}

func getOperationName(c *gin.Context) string {
	return fmt.Sprintf("/%s%s", c.Request.Method, c.FullPath())
}
