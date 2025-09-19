package config

import (
	"log"
	"os"
	"sync"

	"github.com/dhamith93/Skopius/internal/monitor"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Services []monitor.Service `yaml:"services"`
}

var (
	instance *Config
	once     sync.Once
)

func Load(path string) *Config {
	once.Do(func() {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read config file: %v", err)
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("failed to parse config: %v", err)
		}

		instance = &cfg
		log.Println("config loaded")
	})
	return instance
}
