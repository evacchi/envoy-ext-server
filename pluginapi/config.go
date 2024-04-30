package pluginapi

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Config represents the root structure of a config file.
type Config struct {
	Listen string `yaml:"listen"`
	// FilterChains is a collection of sequences of filters.
	FilterChains []FilterChain `yaml:"filter_chains"`
}

// FilterChain is a subsection of Config that defines a single sequence of filters.
type FilterChain struct {
	// Name is the optional name of this FilterChain.
	Name string `yaml:"name"`
	// Filters is a collection of FilterConfig.
	Filters []FilterConfig `yaml:"filters"`
}

// FilterConfig represents the config a single filter.
type FilterConfig struct {
	// Name is the name of the filter.
	Name string `yaml:"name"`
	// Type is the type of filter: it can be `built-in` or `go-plugin`.
	Type string `yaml:"type"`
	// Path is an optional path relative to the config file, pointing to a .so file.
	// It is only valid when Type is `plugin`.
	Path string `yaml:"path"`
	// Config is a custom, unparsed node that is passed over to the plug-in
	// configuration.
	Config map[string]any `yaml:"config"`
}

// ReadConfig reads the configuration from a file.
func ReadConfig(fname string) (*Config, error) {
	c, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &Config{}

	err = yaml.Unmarshal(c, cfg)
	return cfg, err
}
