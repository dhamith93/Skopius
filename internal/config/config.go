package config

import (
	"os"

	"github.com/dhamith93/Skopius/internal/monitor"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Services []monitor.Service `yaml:"services"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
