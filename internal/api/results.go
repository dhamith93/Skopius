package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dhamith93/Skopius/internal/models"
)

// GET /api/v1/results
func ResultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var res models.CheckResult
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		res.Received = time.Now().UTC()

		log.Printf("Received result: %+v", res)
		log.Printf("[%s] %s (code=%d, latency=%dms, err=%s)",
			res.Service, res.Status, res.Code, res.Latency, res.Error)

		// TODO: Store this

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status":"ok"}`))
	}
}
