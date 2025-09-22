package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dhamith93/Skopius/internal/monitor"
	"github.com/google/uuid"
)

type Agent struct {
	ID        string
	Hostname  string
	Region    string
	Token     string
	CreatedAt time.Time
	LastSeen  time.Time
}

type RegisterRequest struct {
	Hostname string `json:"hostname"`
	Region   string `json:"region"`
}

type RegisterResponse struct {
	AgentID string `json:"agent_id"`
	Token   string `json:"token"`
}

type ConfigResponse struct {
	AgentID  string            `json:"agent_id"`
	Services []monitor.Service `json:"checks"`
}

// POST /api/v1/register
func (a *API) RegisterAgentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		agentID := uuid.New().String()
		token := uuid.New().String()

		err := a.Store.RegisterAgent(agentID, req.Hostname, req.Region, token)
		if err != nil {
			http.Error(w, "failed to register agent", http.StatusInternalServerError)
			return
		}

		resp := RegisterResponse{AgentID: agentID, Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// GET /api/v1/config
func (a *API) ConfigHandler(services []monitor.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agentID := r.URL.Query().Get("agent_id")
		token := r.Header.Get("Authorization")

		if agentID == "" || token == "" {
			http.Error(w, "missing agent_id or token", http.StatusBadRequest)
			return
		}

		// strip "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		dbToken, err := a.Store.GetAgentByID(agentID)
		if err != nil {
			http.Error(w, "agent not found", http.StatusUnauthorized)
			return
		}

		if dbToken != token {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		a.Store.UpdateLastSeen(agentID)

		resp := ConfigResponse{
			AgentID:  agentID,
			Services: services,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
