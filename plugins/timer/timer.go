package main

import (
	ep "github.com/evacchi/envoy-ext-server/extproc"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"strconv"
	"time"
)

func New() pluginapi.Plugin {
	return &timerRequestProcessor{}
}

type timerRequestProcessor struct {
	opts *ep.ProcessingOptions
}

func (s *timerRequestProcessor) GetName() string {
	return "timer"
}

func (s *timerRequestProcessor) GetOptions() *ep.ProcessingOptions {
	return s.opts
}

func (s *timerRequestProcessor) ProcessRequestHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	ctx.OverwriteHeader("x-extproc-started-ns", ep.HeaderValue{RawValue: []byte(strconv.FormatInt(ctx.Started.UnixNano(), 10))})

	return ctx.ContinueRequest()
}

func (s *timerRequestProcessor) ProcessRequestBody(ctx *ep.RequestContext, body []byte) error {
	return ctx.ContinueRequest()
}

func (s *timerRequestProcessor) ProcessRequestTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *timerRequestProcessor) ProcessResponseHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	finished := time.Now()
	duration := time.Since(ctx.Started)

	ctx.AddHeader("x-extproc-started-ns", ep.HeaderValue{RawValue: []byte(strconv.FormatInt(ctx.Started.UnixNano(), 10))})
	ctx.AddHeader("x-extproc-finished-ns", ep.HeaderValue{RawValue: []byte(strconv.FormatInt(finished.UnixNano(), 10))})
	ctx.AddHeader("x-upstream-duration-ns", ep.HeaderValue{RawValue: []byte(strconv.FormatInt(duration.Nanoseconds(), 10))})

	return ctx.ContinueRequest()
}

func (s *timerRequestProcessor) ProcessResponseBody(ctx *ep.RequestContext, body []byte) error {
	finished := time.Now()
	duration := time.Since(ctx.Started)

	ctx.OverwriteHeader("x-extproc-started-ns", ep.HeaderValue{RawValue: []byte(strconv.FormatInt(ctx.Started.UnixNano(), 10))})
	ctx.OverwriteHeader("x-extproc-finished-ns", ep.HeaderValue{RawValue: []byte(strconv.FormatInt(finished.UnixNano(), 10))})
	ctx.OverwriteHeader("x-upstream-duration-ns", ep.HeaderValue{RawValue: []byte(strconv.FormatInt(duration.Nanoseconds(), 10))})

	return ctx.ContinueRequest()
}

func (s *timerRequestProcessor) ProcessResponseTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *timerRequestProcessor) Init(opts *ep.ProcessingOptions, nonFlagArgs []string, config pluginapi.FilterConfig) error {
	s.opts = opts
	return nil
}

func (s *timerRequestProcessor) Finish() {}
