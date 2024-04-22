package wasm

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	ep "github.com/evacchi/envoy-ext-server/extproc"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed wasm/hello.wasm
var wasmBin []byte

func NewWasmRequestProcessor() pluginapi.Plugin {
	return &wasmRequestProcessor{}
}

type wasmRequestProcessor struct {
	opts *ep.ProcessingOptions
	rt   wazero.Runtime
	m    wazero.CompiledModule
}

func (s *wasmRequestProcessor) GetName() string {
	return "wasm"
}

func (s *wasmRequestProcessor) GetOptions() *ep.ProcessingOptions {
	return s.opts
}

func (s *wasmRequestProcessor) ProcessRequestHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *wasmRequestProcessor) ProcessRequestBody(ctx *ep.RequestContext, body []byte) error {
	// FIXME: receive from caller?
	c := context.Background()

	stdout := bytes.NewBuffer([]byte{})
	args := []string{"wasm"}
	for k, v := range ctx.AllHeaders.RawHeaders {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}
	cfg := wazero.NewModuleConfig().
		WithStdout(stdout).
		WithStdin(bytes.NewReader(body)).
		WithArgs(args...)

	_, err := s.rt.InstantiateModule(c, s.m, cfg)

	if err != nil {
		return ctx.CancelRequest(500, map[string]ep.HeaderValue{}, "Internal error: "+err.Error())
	}
	err = ctx.ReplaceBodyChunk(stdout.Bytes())
	if err != nil {
		return ctx.CancelRequest(500, map[string]ep.HeaderValue{}, "Internal error: "+err.Error())
	}
	return ctx.ContinueRequest()
}

func (s *wasmRequestProcessor) ProcessRequestTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *wasmRequestProcessor) ProcessResponseHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *wasmRequestProcessor) ProcessResponseBody(ctx *ep.RequestContext, body []byte) error {
	return ctx.ContinueRequest()
}

func (s *wasmRequestProcessor) ProcessResponseTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *wasmRequestProcessor) Init(opts *ep.ProcessingOptions, nonFlagArgs []string) error {
	s.opts = opts
	c := context.Background()
	rt := wazero.NewRuntime(c)
	wasi_snapshot_preview1.MustInstantiate(c, rt)

	m, err := rt.CompileModule(c, wasmBin)
	if err != nil {
		return err
	}

	s.rt = rt
	s.m = m

	return nil
}

func (s *wasmRequestProcessor) Finish() {}
