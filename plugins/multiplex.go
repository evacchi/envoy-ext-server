package plugins

import (
	"github.com/evacchi/envoy-ext-server/extproc"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"log"
)

func NewMultiplexRequestProcessor(names []string, processors []extproc.RequestProcessor) pluginapi.Plugin {
	return &MultiplexRequestProcessor{names: names, processors: processors}
}

type MultiplexRequestProcessor struct {
	names      []string
	processors []extproc.RequestProcessor
}

func (m MultiplexRequestProcessor) Init(opts *extproc.ProcessingOptions, nonFlagArgs []string) error {
	log.Println("Multiplex Init")
	return nil
}

func (m MultiplexRequestProcessor) Finish() {
}

func (m MultiplexRequestProcessor) GetName() string {
	return "multiplex"
}

func (m MultiplexRequestProcessor) GetOptions() *extproc.ProcessingOptions {
	return nil
}

func (m MultiplexRequestProcessor) ProcessRequestHeaders(ctx *extproc.RequestContext, headers extproc.AllHeaders) error {
	for _, p := range m.processors {
		err := p.ProcessRequestHeaders(ctx, headers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MultiplexRequestProcessor) ProcessRequestTrailers(ctx *extproc.RequestContext, trailers extproc.AllHeaders) error {
	for _, p := range m.processors {
		err := p.ProcessRequestTrailers(ctx, trailers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MultiplexRequestProcessor) ProcessResponseHeaders(ctx *extproc.RequestContext, headers extproc.AllHeaders) error {
	for _, p := range m.processors {
		err := p.ProcessResponseHeaders(ctx, headers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MultiplexRequestProcessor) ProcessResponseTrailers(ctx *extproc.RequestContext, trailers extproc.AllHeaders) error {
	for _, p := range m.processors {
		err := p.ProcessResponseTrailers(ctx, trailers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MultiplexRequestProcessor) ProcessResponseBody(ctx *extproc.RequestContext, body []byte) error {
	for _, p := range m.processors {
		err := p.ProcessResponseBody(ctx, body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MultiplexRequestProcessor) ProcessRequestBody(ctx *extproc.RequestContext, body []byte) error {
	for _, p := range m.processors {
		err := p.ProcessRequestBody(ctx, body)
		if err != nil {
			return err
		}
	}
	return nil
}
