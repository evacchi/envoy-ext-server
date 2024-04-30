package wasm

import (
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	config := pluginapi.FilterConfig{
		Name: "wasm",
		Type: "built-in",
		Config: map[string]any{
			"name": "hello",
			"path": "wasm/hello.wasm",
		},
	}

	require.NotPanics(t, func() {
		_ = New(config)
	})
}
