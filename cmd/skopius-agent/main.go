package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dhamith93/Skopius/internal/api"
	"github.com/dhamith93/Skopius/internal/scheduler"
)

const (
	serverURL   = "http://localhost:8080"
	agentConfig = "agent.json"
)

type AgentCredentials struct {
	AgentID string `json:"agent_id"`
	Token   string `json:"token"`
}

func loadCredentials() (*AgentCredentials, error) {
	data, err := os.ReadFile(agentConfig)
	if err != nil {
		return nil, err
	}
	var creds AgentCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}

func saveCredentials(creds *AgentCredentials) error {
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(agentConfig, data, 0644)
}

func registerAgent(hostname, region string) (*AgentCredentials, error) {
	reqBody, _ := json.Marshal(api.RegisterRequest{Hostname: hostname, Region: region})
	resp, err := http.Post(serverURL+"/api/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registration failed: %s", resp.Status)
	}

	var regResp api.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		return nil, err
	}

	return &AgentCredentials{AgentID: regResp.AgentID, Token: regResp.Token}, nil
}

func fetchConfig(creds *AgentCredentials) (*api.ConfigResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/config?agent_id=%s", serverURL, creds.AgentID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+creds.Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("config fetch failed: %s", resp.Status)
	}

	var cfg api.ConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {
	log.Println("Starting Skopius agent...")
	creds, err := loadCredentials()
	if err != nil {
		log.Println("No local credentials, registering new agent...")

		creds, err = registerAgent("agent-1", "us-east-1")
		if err != nil {
			log.Fatal("Failed to register:", err)
		}
		if err := saveCredentials(creds); err != nil {
			log.Fatal("Failed to save credentials:", err)
		}
		log.Println("Agent registered and credentials saved.")
	} else {
		log.Println("Loaded existing credentials:", creds.AgentID)
	}

	cfg, err := fetchConfig(creds)
	if err != nil {
		log.Fatal("Error fetching config:", err)
	}
	log.Printf("Got config for agent %s: %+v\n", cfg.AgentID, cfg.Services)

	scheduler := scheduler.NewScheduler(cfg.Services)
	go scheduler.Start()

	go func() {
		for res := range scheduler.Results {
			log.Printf("[%s] %s (code=%d, latency=%dms, err=%s)",
				res.Name, res.Status, res.Code, res.Latency, res.Error)
			// handle results
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Skopius...")
	scheduler.Stop()
}
