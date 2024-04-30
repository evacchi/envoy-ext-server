package testing

import "github.com/evacchi/envoy-ext-server/pluginapi"

type Context struct {
	Scheme    string
	Authority string
	Method    string
	Path      string
	FullPath  string
	RequestID string
	Headers   pluginapi.Headers
	Trailers  pluginapi.Trailers
	Body      pluginapi.Body
}

type RequestContext struct {
	Context
}

func (c *RequestContext) Scheme() string               { return c.Context.Scheme }
func (c *RequestContext) Authority() string            { return c.Context.Authority }
func (c *RequestContext) Method() string               { return c.Context.Method }
func (c *RequestContext) Path() string                 { return c.Context.Path }
func (c *RequestContext) FullPath() string             { return c.Context.FullPath }
func (c *RequestContext) RequestID() string            { return c.Context.RequestID }
func (c *RequestContext) Headers() pluginapi.Headers   { return c.Context.Headers }
func (c *RequestContext) Trailers() pluginapi.Trailers { return c.Context.Trailers }
func (c *RequestContext) Body() pluginapi.Body         { return c.Context.Body }
func (c *RequestContext) IsRequest() bool              { return true }

type Headers map[string][]byte

func (h Headers) Get(key string) pluginapi.Header {
	return pluginapi.Header{Key: key, Value: h[key]}
}

func (h Headers) SetRaw(key string, bytes []byte) {
	h[key] = bytes
}

func (h Headers) All() []pluginapi.Header {
	var headers []pluginapi.Header
	for k := range h {
		headers = append(headers, h.Get(k))
	}
	return headers
}

type Trailers map[string]string

type Body struct {
	Content []byte
}

func (b *Body) Get() []byte {
	return b.Content
}

func (b *Body) Set(data []byte) {
	b.Content = data
}
