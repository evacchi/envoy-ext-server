package main

import (
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins/timer"
)

func New(config pluginapi.FilterConfig) pluginapi.Plugin {
	return timer.New(config)
}
