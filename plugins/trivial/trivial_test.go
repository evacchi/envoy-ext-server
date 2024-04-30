package trivial

import (
	"github.com/evacchi/envoy-ext-server/pluginapi"
	testingapi "github.com/evacchi/envoy-ext-server/pluginapi/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	require.NotPanics(t, func() {
		_ = New(pluginapi.FilterConfig{})
	})
}

func TestTrivial_OnRequestBody(t *testing.T) {
	tests := []struct {
		name        string
		req         pluginapi.RequestContext
		expectedReq pluginapi.RequestContext
	}{
		{
			name: "empty request context",
			req: &testingapi.RequestContext{
				Context: testingapi.Context{
					Scheme:    "http",
					Authority: "localhost",
					Method:    "POST",
					Path:      "/",
					FullPath:  "/",
					RequestID: "123",
					Headers: testingapi.Headers{
						"x-trivial-process": []byte("yes"),
					},
					Body: &testingapi.Body{
						Content: []byte(`{"key":"value"}`),
					},
				},
			},
			expectedReq: &testingapi.RequestContext{
				Context: testingapi.Context{
					Scheme:    "http",
					Authority: "localhost",
					Method:    "POST",
					Path:      "/",
					FullPath:  "/",
					RequestID: "123",
					Headers: testingapi.Headers{
						"x-trivial-process": []byte("yes"),
					},
					Body: &testingapi.Body{
						Content: []byte(`{"key":"value","message":"hello from trivial plugin"}`),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			p := New(pluginapi.FilterConfig{})
			req := tc.req
			err := p.OnRequestBody(req)
			require.NoError(t, err)
			require.Equal(t, tc.expectedReq, req)
		})
	}
}
