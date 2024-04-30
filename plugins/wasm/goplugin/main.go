package main

import (
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins/wasm"
)

func New(config pluginapi.FilterConfig) pluginapi.Plugin {
	return wasm.New(config)
}
