package extproc

import (
	"fmt"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	extprocv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"slices"
)

type reqContext struct {
	rContext
}

func (r *reqContext) IsRequest() bool { return true }

var _ pluginapi.RequestContext = (*reqContext)(nil)

func newReqContext(headers *corev3.HeaderMap) (*reqContext, error) {
	context, err := newRContext(headers)
	if err != nil {
		return nil, err
	}
	return &reqContext{*context}, nil
}

func (r *reqContext) toResponseRequestHeaders() *extprocv3.ProcessingResponse {
	aa := corev3.HeaderValueOption_HeaderAppendAction(
		corev3.HeaderValueOption_HeaderAppendAction_value["OVERWRITE_IF_EXISTS_OR_ADD"],
	)

	h, b := r.headers, r.body

	var headers []*corev3.HeaderValueOption
	for i, v := range h.values {
		if slices.Contains(h.updated, i) {
			continue
		}
		headers = append(headers, &corev3.HeaderValueOption{
			Header: &corev3.HeaderValue{
				Key:      v.Key,
				RawValue: v.Value,
			},
			AppendAction: aa,
		})
	}

	if b.updated {
		headers = append(headers, &corev3.HeaderValueOption{
			Header: &corev3.HeaderValue{
				Key:      kContentLength,
				RawValue: []byte(fmt.Sprintf("%d", len(b.body))),
			}})
	}
	b.updated = false

	return &extprocv3.ProcessingResponse{
		Response: &extprocv3.ProcessingResponse_RequestHeaders{
			RequestHeaders: &extprocv3.HeadersResponse{
				Response: &extprocv3.CommonResponse{
					HeaderMutation: &extprocv3.HeaderMutation{
						SetHeaders: headers,
					},
					BodyMutation: &extprocv3.BodyMutation{
						Mutation: &extprocv3.BodyMutation_Body{
							Body: b.body,
						},
					},
				},
			},
		},
	}
}

func (r *reqContext) toResponseRequestTrailers() *extprocv3.ProcessingResponse {
	t := r.trailers

	aa := corev3.HeaderValueOption_HeaderAppendAction(
		corev3.HeaderValueOption_HeaderAppendAction_value["OVERWRITE_IF_EXISTS_OR_ADD"],
	)

	var headers []*corev3.HeaderValueOption
	for i, v := range t.values {
		if !t.isUpdated(i) {
			continue
		}
		headers = append(headers, &corev3.HeaderValueOption{
			Header: &corev3.HeaderValue{
				Key:      v.Key,
				RawValue: v.Value,
			},
			AppendAction: aa,
		})
	}

	return &extprocv3.ProcessingResponse{
		Response: &extprocv3.ProcessingResponse_RequestTrailers{
			RequestTrailers: &extprocv3.TrailersResponse{
				HeaderMutation: &extprocv3.HeaderMutation{
					SetHeaders: headers,
				},
			},
		},
	}
}

func (r *reqContext) toResponseRequestBody() *extprocv3.ProcessingResponse {
	aa := corev3.HeaderValueOption_HeaderAppendAction(
		corev3.HeaderValueOption_HeaderAppendAction_value["OVERWRITE_IF_EXISTS_OR_ADD"],
	)

	h, b := r.headers, r.body

	var headers []*corev3.HeaderValueOption
	for i, k := range h.values {
		if !h.isUpdated(i) {
			continue
		}
		headers = append(headers, &corev3.HeaderValueOption{
			Header: &corev3.HeaderValue{
				Key:      k.Key,
				RawValue: k.Value,
			},
			AppendAction: aa,
		})
	}

	if b.updated {
		headers = append(headers, &corev3.HeaderValueOption{
			Header: &corev3.HeaderValue{
				Key:      kContentLength,
				RawValue: []byte(fmt.Sprintf("%d", len(b.body))),
			}})
	}
	b.updated = false

	return &extprocv3.ProcessingResponse{
		Response: &extprocv3.ProcessingResponse_RequestBody{
			RequestBody: &extprocv3.BodyResponse{
				Response: &extprocv3.CommonResponse{
					HeaderMutation: &extprocv3.HeaderMutation{
						SetHeaders: headers,
					},
					BodyMutation: &extprocv3.BodyMutation{
						Mutation: &extprocv3.BodyMutation_Body{
							Body: b.body,
						},
					},
				},
			},
		},
	}
}
