package noop

import (
	"github.com/evacchi/envoy-ext-server/pluginapi"
)

type noop struct {
	pluginapi.DefaultPlugin
}

func New(config pluginapi.FilterConfig) pluginapi.Plugin {
	return &noop{}
}
