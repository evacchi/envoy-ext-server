package plugins

import (
	"github.com/evacchi/envoy-ext-server/extproc"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"log"
)

func NewFilterChain(processors []pluginapi.Plugin) pluginapi.Plugin {
	return &FilterChain{processors: processors}
}

type FilterChain struct {
	processors []pluginapi.Plugin
}

func (m FilterChain) Init(opts *extproc.ProcessingOptions, nonFlagArgs []string) error {
	log.Println("FilterChain Init")
	for _, p := range m.processors {
		if err := p.Init(opts, nonFlagArgs); err != nil {
			return err
		}
	}
	return nil
}

func (m FilterChain) Finish() {
	for _, p := range m.processors {
		p.Finish()
	}
}

func (m FilterChain) GetName() string {
	return "multiplex"
}

func (m FilterChain) GetOptions() *extproc.ProcessingOptions {
	return nil
}

func (m FilterChain) ProcessRequestHeaders(ctx *extproc.RequestContext, headers extproc.AllHeaders) error {
	for _, p := range m.processors {
		err := p.ProcessRequestHeaders(ctx, headers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m FilterChain) ProcessRequestTrailers(ctx *extproc.RequestContext, trailers extproc.AllHeaders) error {
	for _, p := range m.processors {
		err := p.ProcessRequestTrailers(ctx, trailers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m FilterChain) ProcessResponseHeaders(ctx *extproc.RequestContext, headers extproc.AllHeaders) error {
	for _, p := range m.processors {
		err := p.ProcessResponseHeaders(ctx, headers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m FilterChain) ProcessResponseTrailers(ctx *extproc.RequestContext, trailers extproc.AllHeaders) error {
	for _, p := range m.processors {
		err := p.ProcessResponseTrailers(ctx, trailers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m FilterChain) ProcessResponseBody(ctx *extproc.RequestContext, body []byte) error {
	for _, p := range m.processors {
		err := p.ProcessResponseBody(ctx, body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m FilterChain) ProcessRequestBody(ctx *extproc.RequestContext, body []byte) error {
	for _, p := range m.processors {
		err := p.ProcessRequestBody(ctx, body)
		if err != nil {
			return err
		}
	}
	return nil
}
