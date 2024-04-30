package pluginapi

import (
	"fmt"
	"plugin"
)

type PluginFactory = func(config FilterConfig) Plugin

func FromSharedObject(fname string) (PluginFactory, error) {
	p, err := plugin.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %v", fname, err)
	}
	ff, err := p.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup New function for plugin %s: %v", err, fname)
	}
	pf := ff.(PluginFactory)
	return pf, nil
}

type Header struct {
	Key   string
	Value []byte
}

type Headers interface {
	Get(key string) Header
	SetRaw(key string, bytes []byte)
	All() []Header
}

type Trailers interface {
	Get(key string) Header
	SetRaw(key string, bytes []byte)
	All() []Header
}

type Body interface {
	Get() []byte
	Set(data []byte)
}

type RequestContext interface {
	Scheme() string
	Authority() string
	Method() string
	Path() string
	FullPath() string
	RequestID() string
	Headers() Headers
	Trailers() Trailers
	Body() Body
	IsRequest() bool
}

type RequestHandler interface {
	OnRequestHeaders(req RequestContext) error
	OnRequestBody(req RequestContext) error
	OnRequestTrailers(req RequestContext) error
}

type NoopRequestHandler struct{}

func (u NoopRequestHandler) OnRequestHeaders(req RequestContext) error  { return nil }
func (u NoopRequestHandler) OnRequestBody(req RequestContext) error     { return nil }
func (u NoopRequestHandler) OnRequestTrailers(req RequestContext) error { return nil }

type ResponseContext interface {
	Scheme() string
	Authority() string
	Method() string
	Path() string
	FullPath() string
	RequestID() string
	Headers() Headers
	Trailers() Trailers
	Body() Body
	IsResponse() bool
}

type ResponseHandler interface {
	OnResponseHeaders(req ResponseContext) error
	OnResponseBody(req ResponseContext) error
	OnResponseTrailers(req ResponseContext) error
}

type NoopResponseHandler struct{}

func (u NoopResponseHandler) OnResponseHeaders(req ResponseContext) error  { return nil }
func (u NoopResponseHandler) OnResponseBody(req ResponseContext) error     { return nil }
func (u NoopResponseHandler) OnResponseTrailers(req ResponseContext) error { return nil }

type Plugin interface {
	RequestHandler
	ResponseHandler
}

type DefaultPlugin struct {
	NoopRequestHandler
	NoopResponseHandler
}
