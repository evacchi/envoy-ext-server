package pluginapi

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	FilterChains []Filter `yaml:"filter_chains"`
}

type Filter struct {
	Filters []FilterConfig `yaml:"filters"`
}

type FilterConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

func ReadConfig(fname string) (*Config, error) {
	c, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &Config{}

	err = yaml.Unmarshal(c, cfg)
	return cfg, err
}
