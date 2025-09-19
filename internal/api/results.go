package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type CheckResult struct {
	AgentID string    `json:"agent_id"`
	Service string    `json:"service"`
	URL     string    `json:"url"`
	Status  string    `json:"status"`
	Code    int       `json:"code"`
	Latency int64     `json:"latency"` // ms
	Error   string    `json:"error"`
	Time    time.Time `json:"time"`
}

// GET /api/v1/results
func ResultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var res CheckResult
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		res.Time = time.Now().UTC()

		log.Printf("Received result: %+v", res)
		log.Printf("[%s] %s (code=%d, latency=%dms, err=%s)",
			res.Service, res.Status, res.Code, res.Latency, res.Error)

		// TODO: Store this

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status":"ok"}`))
	}
}
