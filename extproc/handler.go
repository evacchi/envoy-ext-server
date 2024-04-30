package extproc

import (
	"fmt"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"slices"
	"strings"
)

const kContentLength = "Content-Length"

type headerTrailer struct {
	keyMap  map[string]int
	values  []pluginapi.Header
	updated []int
}

type headers struct {
	headerTrailer
}
type trailers struct {
	headerTrailer
}
type body struct {
	body    []byte
	updated bool
}

func (h *headerTrailer) Get(key string) pluginapi.Header {
	return h.values[h.keyMap[key]]
}

func (h *headerTrailer) SetRaw(key string, bytes []byte) {
	h.setRaw(pluginapi.Header{Key: key, Value: bytes})
	idx := h.keyMap[key]
	h.setUpdated(idx)
}

func (h *headerTrailer) setRaw(header pluginapi.Header) {
	h.values = append(h.values, header)
	h.keyMap[header.Key] = len(h.values) - 1
}

func (h *headerTrailer) All() []pluginapi.Header {
	vs := make([]pluginapi.Header, len(h.values))
	for i := range h.values {
		vs[i] = h.values[i]
	}
	return vs
}

func (h *headerTrailer) setUpdated(idx int) {
	if h.isUpdated(idx) {
		return
	}
	h.updated = append(h.updated, idx)
}

func (h *headerTrailer) isUpdated(idx int) bool {
	return slices.Contains(h.updated, idx)
}

func (b *body) Get() []byte {
	return b.body
}

func (b *body) Set(data []byte) {
	b.updated = true
	b.body = data
}

type rContext struct {
	scheme    string
	authority string
	method    string
	path      string
	fullPath  string
	requestID string

	headers  *headers
	trailers *trailers
	body     *body
}

func (r *rContext) Scheme() string               { return r.scheme }
func (r *rContext) Authority() string            { return r.authority }
func (r *rContext) Method() string               { return r.method }
func (r *rContext) Path() string                 { return r.path }
func (r *rContext) FullPath() string             { return r.fullPath }
func (r *rContext) RequestID() string            { return r.requestID }
func (r *rContext) Headers() pluginapi.Headers   { return r.headers }
func (r *rContext) Trailers() pluginapi.Trailers { return r.trailers }
func (r *rContext) Body() pluginapi.Body         { return r.body }

func newRContext(headers *corev3.HeaderMap) (*rContext, error) {
	eitherValue := func(h *corev3.HeaderValue) string {
		if h == nil {
			return ""
		}
		val := h.Value
		if len(h.RawValue) > 0 {
			val = string(h.RawValue)
		}
		return val
	}

	var err error
	r := &rContext{}
	r.body = &body{}
	r.headers, err = toHeaders(headers)
	if err != nil {
		return nil, fmt.Errorf("parse header is failed: %w", err)
	}

	for _, h := range headers.Headers {
		switch h.Key {
		case ":scheme":
			r.scheme = eitherValue(h)

		case ":authority":
			r.authority = eitherValue(h)

		case ":method":
			r.method = eitherValue(h)

		case ":path":
			r.fullPath = eitherValue(h)
			r.path = strings.Split(r.fullPath, "?")[0]

		case "x-request-id":
			r.requestID = eitherValue(h)

		default:
		}
	}

	return r, nil
}

func toHeaderTrailer(headerMap *corev3.HeaderMap) (headers *headerTrailer, err error) {
	headers = &headerTrailer{keyMap: make(map[string]int)}

	for _, h := range headerMap.Headers {
		if len(h.Value) > 0 && len(h.RawValue) > 0 {
			err = fmt.Errorf("only one of 'value' or 'raw_value' can be set")
			return
		}

		if len(h.Value) > 0 {
			panic("unsupported string value in HeaderMap	")
		}

		headers.setRaw(pluginapi.Header{Key: h.Key, Value: h.RawValue})
	}
	return
}

func toHeaders(headerMap *corev3.HeaderMap) (*headers, error) {
	t, err := toHeaderTrailer(headerMap)
	if err != nil {
		return nil, err
	}
	return &headers{*t}, nil
}

func toTrailers(headerMap *corev3.HeaderMap) (*trailers, error) {
	t, err := toHeaderTrailer(headerMap)
	if err != nil {
		return nil, err
	}
	return &trailers{*t}, nil
}
