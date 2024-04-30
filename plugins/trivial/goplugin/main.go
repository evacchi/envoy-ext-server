package main

import (
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins/trivial"
)

func New(config pluginapi.FilterConfig) pluginapi.Plugin {
	return trivial.New(config)
}
