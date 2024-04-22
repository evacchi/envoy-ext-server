package main

import (
	"crypto/sha256"
	"encoding/hex"
	ep "github.com/evacchi/envoy-ext-server/extproc"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"hash"
)

func New() pluginapi.Plugin {
	return &digestRequestProcessor{}
}

type digestRequestProcessor struct {
	opts *ep.ProcessingOptions
}

func GetHasher(ctx *ep.RequestContext) (hash.Hash, error) {
	val, err := ctx.GetValue("hasher")
	if err != nil {
		return nil, err
	}
	return val.(hash.Hash), nil
}

func GetDigest(ctx *ep.RequestContext) (string, error) {
	val, err := ctx.GetValue("digest")
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func (s *digestRequestProcessor) GetName() string {
	return "digest"
}

func (s *digestRequestProcessor) GetOptions() *ep.ProcessingOptions {
	return s.opts
}

func (s *digestRequestProcessor) ProcessRequestHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	hasher := sha256.New()
	ctx.SetValue("hasher", hasher)

	hasher.Write([]byte(ctx.Method + ":" + ctx.Path)) // method:path

	if ctx.EndOfStream {
		digest := hex.EncodeToString(hasher.Sum(nil))
		ctx.SetValue("digest", digest)
		ctx.AddHeader("x-extproc-request-digest", ep.HeaderValue{RawValue: []byte(digest)})
	}

	return ctx.ContinueRequest()
}

func (s *digestRequestProcessor) ProcessRequestBody(ctx *ep.RequestContext, body []byte) error {
	hasher, _ := GetHasher(ctx)
	hasher.Write([]byte(":"))
	hasher.Write(body)

	if ctx.EndOfStream {
		digest := hex.EncodeToString(hasher.Sum(nil))
		ctx.SetValue("digest", digest)
		ctx.AddHeader("x-extproc-request-digest", ep.HeaderValue{RawValue: []byte(digest)})
	}
	return ctx.ContinueRequest()
}

func (s *digestRequestProcessor) ProcessRequestTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *digestRequestProcessor) ProcessResponseHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	if ctx.EndOfStream {
		digest, _ := GetDigest(ctx)
		ctx.AddHeader("x-extproc-request-digest", ep.HeaderValue{RawValue: []byte(digest)})
	}
	return ctx.ContinueRequest()
}

func (s *digestRequestProcessor) ProcessResponseBody(ctx *ep.RequestContext, body []byte) error {
	if ctx.EndOfStream {
		digest, _ := GetDigest(ctx)
		ctx.AddHeader("x-extproc-request-digest", ep.HeaderValue{RawValue: []byte(digest)})
	}
	return ctx.ContinueRequest()
}

func (s *digestRequestProcessor) ProcessResponseTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *digestRequestProcessor) Init(opts *ep.ProcessingOptions, nonFlagArgs []string, config pluginapi.FilterConfig) error {
	s.opts = opts
	return nil
}

func (s *digestRequestProcessor) Finish() {}
