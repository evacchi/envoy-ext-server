package pluginapi

type PluginFactory func() Plugin
type Registry map[string]PluginFactory

var registry = Registry{}

func Register(name string, mk PluginFactory) {
	registry[name] = mk
}

func Instantiate(name string, config FilterConfig) Plugin {
	if p, ok := registry[name]; ok {
		return p()
	}
	return nil
}

func InstantiateAll() (instances []Plugin) {
	for _, p := range registry {
		instances = append(instances, p())
	}
	return instances
}
