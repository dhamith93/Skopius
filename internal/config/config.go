package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Service struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Interval int    `yaml:"interval"` // seconds
}

type Config struct {
	Services []Service `yaml:"services"`
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
