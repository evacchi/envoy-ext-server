package plugins

import (
	"fmt"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins/timer"
	"github.com/evacchi/envoy-ext-server/plugins/trivial"
	"github.com/evacchi/envoy-ext-server/plugins/wasm"
)

var builtIns = map[string]pluginapi.PluginFactory{
	"trivial": trivial.New,
	"timer":   timer.New,
	"wasm":    wasm.New,
}

// NewFilterChain creates a new FilterChain with the given configurations, for the given pluginFactories.
// Each index in the config corresponds to the config of the plugin factory at the same index.
func NewFilterChain(config *pluginapi.FilterChain, pluginFactories []pluginapi.PluginFactory) (pluginapi.Plugin, error) {
	if len(config.Filters) != len(pluginFactories) {
		panic(fmt.Errorf("number of configured filters does not match number of plugin factories"))
	}
	fc := FilterChain{
		pluginFactories: pluginFactories,
		plugins:         make([]pluginapi.Plugin, len(pluginFactories)),
		config:          config,
	}
	err := fc.Reload()
	return &fc, err
}

// NewFilterChainFromConfig creates a new FilterChain, loading the configuration from the given path.
func NewFilterChainFromConfig(config *pluginapi.Config) (pluginapi.Plugin, error) {
	ps, configs, err := loadFilterChain(config)
	if err != nil {
		return nil, err
	}
	return NewFilterChain(configs, ps)
}

func loadFilterChain(config *pluginapi.Config) ([]pluginapi.PluginFactory, *pluginapi.FilterChain, error) {
	var pfs []pluginapi.PluginFactory
	for _, fc := range config.FilterChains {
		for _, f := range fc.Filters {
			if pf, err := loadPlugin(f, builtIns); err != nil {
				return nil, nil, err
			} else {
				pfs = append(pfs, pf)
			}
		}
		return pfs, &fc, nil
	}
	return nil, nil, fmt.Errorf("no plugin found")
}

func loadPlugin(f pluginapi.FilterConfig, builtIns map[string]pluginapi.PluginFactory) (pluginapi.PluginFactory, error) {
	switch f.Type {
	case "built-in":
		if pf, ok := builtIns[f.Name]; !ok {
			return nil, fmt.Errorf("no such built-in filter: %s", f.Name)
		} else {
			return pf, nil
		}
	case "go-plugin":
		return pluginapi.FromSharedObject(f.Path)
	default:
		return nil, fmt.Errorf("unsupported filter type: %s", f.Type)
	}
}

// FilterChain is a pluginapi.Plugin that invokes all configured filters in sequence.
type FilterChain struct {
	config          *pluginapi.FilterChain
	pluginFactories []pluginapi.PluginFactory
	plugins         []pluginapi.Plugin
}

// Reload is a lifecycle method to reload all the configured plugins.
func (f *FilterChain) Reload() error {
	for i, pf := range f.pluginFactories {
		fc := f.config.Filters[i]
		f.plugins[i] = pf(fc)
	}
	return nil
}

func (f *FilterChain) OnRequestHeaders(req pluginapi.RequestContext) error {
	for _, plugin := range f.plugins {
		if err := plugin.OnRequestHeaders(req); err != nil {
			return err
		}
	}
	return nil
}

func (f *FilterChain) OnRequestBody(req pluginapi.RequestContext) error {
	for _, plugin := range f.plugins {
		if err := plugin.OnRequestBody(req); err != nil {
			return err
		}
	}
	return nil
}

func (f *FilterChain) OnRequestTrailers(req pluginapi.RequestContext) error {
	for _, plugin := range f.plugins {
		if err := plugin.OnRequestTrailers(req); err != nil {
			return err
		}
	}
	return nil
}

func (f *FilterChain) OnResponseHeaders(resp pluginapi.ResponseContext) error {
	for _, plugin := range f.plugins {
		if err := plugin.OnResponseHeaders(resp); err != nil {
			return err
		}
	}
	return nil
}

func (f *FilterChain) OnResponseBody(resp pluginapi.ResponseContext) error {
	for _, plugin := range f.plugins {
		if err := plugin.OnResponseBody(resp); err != nil {
			return err
		}
	}
	return nil
}

func (f *FilterChain) OnResponseTrailers(resp pluginapi.ResponseContext) error {
	for _, plugin := range f.plugins {
		if err := plugin.OnResponseTrailers(resp); err != nil {
			return err
		}
	}
	return nil
}
