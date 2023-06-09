// from https://github.com/SkyAPM/go2sky/blob/master/plugins/http/client.go

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	// "github.com/SkyAPM/go2sky/internal/tool"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const componentIDGOHttpClient = 5005

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	errInvalidTracer            = Error("invalid tracer")
	TagPayload       go2sky.Tag = "payload"
	TagResponse      go2sky.Tag = "response"
)

type ClientConfig struct {
	name      string
	client    *http.Client
	tracer    *go2sky.Tracer
	extraTags map[string]string
}

// ClientOption allows optional configuration of Client.
type ClientOption func(*ClientConfig)

// WithOperationName override default operation name.
func WithClientOperationName(name string) ClientOption {
	return func(c *ClientConfig) {
		c.name = name
	}
}

// WithClientTag adds extra tag to client spans.
func WithClientTag(key string, value string) ClientOption {
	return func(c *ClientConfig) {
		if c.extraTags == nil {
			c.extraTags = make(map[string]string)
		}
		c.extraTags[key] = value
	}
}

// WithClient set customer http client.
func WithClient(client *http.Client) ClientOption {
	return func(c *ClientConfig) {
		c.client = client
	}
}

// NewClient returns an HTTP Client with tracer
func NewClient(tracer *go2sky.Tracer, options ...ClientOption) (*http.Client, error) {
	// if tracer == nil {
	// 	return nil, errInvalidTracer
	// }
	co := &ClientConfig{tracer: tracer}
	for _, option := range options {
		option(co)
	}
	if co.client == nil {
		co.client = &http.Client{}
	}
	tp := &transport{
		ClientConfig: co,
		delegated:    http.DefaultTransport,
	}
	if co.client.Transport != nil {
		tp.delegated = co.client.Transport
	}
	co.client.Transport = tp
	return co.client, nil
}

type transport struct {
	*ClientConfig
	delegated http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	if t.tracer == nil {
		return t.delegated.RoundTrip(req)
	}
	span, err := t.tracer.CreateExitSpan(req.Context(), getOperationName(t.name, req), req.Host, func(key, value string) error {
		req.Header.Set(key, value)
		return nil
	})
	if err != nil {
		return t.delegated.RoundTrip(req)
	}
	defer span.End()
	span.SetComponent(componentIDGOHttpClient)
	for k, v := range t.extraTags {
		span.Tag(go2sky.Tag(k), v)
	}
	span.Tag(go2sky.TagHTTPMethod, req.Method)
	span.Tag(go2sky.TagURL, req.URL.String())

	// tag req payload
	var bodyBytes []byte
	if req.ContentLength > 0 {
		contentType := req.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			bodyBytes, _ = ioutil.ReadAll(req.Body)
			// 新建缓冲区并替换原有Request.body
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		} else {
			r := make(map[string]interface{})
			req.ParseForm()
			for k, v := range req.Form {
				r[k] = v
			}
			bodyBytes, _ = json.Marshal(r)
		}
		span.Tag(TagPayload, string(bodyBytes))
	}

	span.SetSpanLayer(agentv3.SpanLayer_Http)
	res, err = t.delegated.RoundTrip(req)
	if err != nil {
		span.Error(time.Now(), err.Error())
		return
	}
	span.Tag(go2sky.TagStatusCode, strconv.Itoa(res.StatusCode))

	// tag response
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if res.StatusCode >= http.StatusBadRequest {
		span.Error(time.Now(), "Errors on handling client", string(bodyBytes))
	} else {
		span.Tag(TagResponse, string(bodyBytes))
	}

	return res, nil
}

func getOperationName(name string, r *http.Request) string {
	if name == "" {
		return fmt.Sprintf("/%s%s", r.Method, r.URL.Path)
	}
	return name
}
