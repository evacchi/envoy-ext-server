package main

import (
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins/noop"
)

func New(config pluginapi.FilterConfig) pluginapi.Plugin {
	return noop.New(config)
}
