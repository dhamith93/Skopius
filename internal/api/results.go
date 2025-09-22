package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dhamith93/Skopius/internal/models"
	"github.com/dhamith93/Skopius/internal/store"
)

type API struct {
	Store *store.Store
}

// POST /api/v1/results
func (a *API) ResultHandler() http.HandlerFunc {
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

		if err := a.Store.SaveResult(res); err != nil {
			http.Error(w, "failed to save result", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status":"ok"}`))
	}
}

func (a *API) GetLatestResultsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// default limit = 20
		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		results, err := a.Store.GetLatestResults(limit)
		if err != nil {
			http.Error(w, "failed to fetch results", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
