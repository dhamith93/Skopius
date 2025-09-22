package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dhamith93/Skopius/internal/api"
	"github.com/dhamith93/Skopius/internal/config"
	"github.com/dhamith93/Skopius/internal/store"
)

func main() {
	cfg := config.Load("config.yml")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store, err := store.NewStore("skopius.db")
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	api := &api.API{
		Store: store,
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/register", api.RegisterAgentHandler())
	mux.Handle("/api/v1/config", api.ConfigHandler(cfg.Services))
	mux.Handle("/api/v1/results", api.ResultHandler())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Println("Starting Skopius server on :" + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Skopius server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Skopius server forced to shutdown: %v", err)
	}

}
