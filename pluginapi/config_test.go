package pluginapi

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name     string
		fname    string
		expected *Config
	}{
		{
			name:  "example.yml",
			fname: "testdata/example.yml",
			expected: &Config{
				FilterChains: []FilterChain{
					{
						Name: "first chain",
						Filters: []FilterConfig{
							{
								Name: "timer",
								Type: "built-in",
							},
							{
								Name: "trivial",
								Type: "built-in",
							},
							{
								Name: "wasm",
								Type: "built-in",
								Config: map[string]any{
									"name": "hello",
									"path": "plugins/wasm/wasm/hello.wasm",
								},
							},
						},
					},
					{
						Name: "second chain",
						Filters: []FilterConfig{
							{
								Name: "timer",
								Type: "go-plugin",
								Path: "plugins/timer/timer.so",
							},
							{
								Name: "trivial",
								Type: "go-plugin",
								Path: "plugins/trivial/trivial.so",
							},
							{
								Name: "wasm",
								Type: "go-plugin",
								Path: "plugins/wasm/wasm.so",
								Config: map[string]any{
									"name": "hello",
									"path": "plugins/wasm/wasm/hello.wasm",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ReadConfig(tt.fname)
			require.NoError(t, err)
			require.Equal(t, tt.expected, cfg)
		})
	}
}
