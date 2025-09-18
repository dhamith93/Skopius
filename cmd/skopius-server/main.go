package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dhamith93/Skopius/internal/config"
	"github.com/dhamith93/Skopius/internal/monitor"
	"github.com/dhamith93/Skopius/internal/scheduler"
)

func main() {
	cfg, err := config.Load("config.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	s := scheduler.New()

	for _, svc := range cfg.Services {
		service := svc
		s.Every(service.Interval, func() {
			result := monitor.CheckHTTP(service.URL)
			fmt.Printf("[%s] %s -> %s (%.2fms)\n",
				time.Now().Format(time.RFC3339),
				service.Name,
				result.Status,
				result.Latency,
			)
		})
	}

	s.Start()
}
