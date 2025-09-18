package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dhamith93/Skopius/internal/config"
	"github.com/dhamith93/Skopius/internal/scheduler"
)

func main() {
	cfg, err := config.Load("config.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("Starting Skopius...")

	scheduler := scheduler.NewScheduler(cfg.Services)
	go scheduler.Start()

	go func() {
		for res := range scheduler.Results {
			log.Printf("[%s] %s (code=%d, latency=%dms, err=%s)",
				res.Name, res.Status, res.Code, res.Latency, res.Error)
			// handle results (store in DB, alerts, etc.)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Skopius...")
	scheduler.Stop()
}
