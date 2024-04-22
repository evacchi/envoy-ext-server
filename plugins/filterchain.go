package plugins

import (
	"github.com/evacchi/envoy-ext-server/extproc"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins/data"
	"github.com/evacchi/envoy-ext-server/plugins/dedup"
	"github.com/evacchi/envoy-ext-server/plugins/digest"
	"github.com/evacchi/envoy-ext-server/plugins/echo"
	"github.com/evacchi/envoy-ext-server/plugins/masker"
	"github.com/evacchi/envoy-ext-server/plugins/noop"
	"github.com/evacchi/envoy-ext-server/plugins/timer"
	"github.com/evacchi/envoy-ext-server/plugins/trivial"
	"github.com/evacchi/envoy-ext-server/plugins/wasm"
	"log"
)

func init() {
	pluginapi.Register("data", data.NewDataRequestProcessor)
	pluginapi.Register("dedup", dedup.NewDedupRequestProcessor)
	pluginapi.Register("digest", digest.NewDigestRequestProcessor)
	pluginapi.Register("echo", echo.NewEchoRequestProcessor)
	pluginapi.Register("masker", masker.NewMaskerRequestProcessor)
	pluginapi.Register("noop", noop.NewNoopRequestProcessor)
	pluginapi.Register("timer", timer.NewTimerRequestProcessor)
	pluginapi.Register("trivial", trivial.NewTrivialRequestProcessor)
	pluginapi.Register("wasm", wasm.NewWasmRequestProcessor)
}

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
