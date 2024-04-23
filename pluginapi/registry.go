package pluginapi

type PluginFactory func() Plugin
type Registry map[string]PluginFactory

var registry = Registry{}

func Register(name string, mk PluginFactory) {
	registry[name] = mk
}

func Instantiate(fc FilterConfig) Plugin {
	if p, ok := registry[fc.Name]; ok {
		plugin := p()
		plugin.Init(nil, nil, fc)
		return plugin
	}
	return nil
}

func InstantiateAll() (instances []Plugin) {
	for _, p := range registry {
		instances = append(instances, p())
	}
	return instances
}
