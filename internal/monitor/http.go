package monitor

import (
	"net/http"
	"time"
)

type Result struct {
	Status  string
	Latency float64
}

func CheckHTTP(url string) Result {
	start := time.Now()
	resp, err := http.Get(url)
	latency := time.Since(start).Seconds() * 1000 // ms

	if err != nil {
		return Result{Status: "DOWN", Latency: latency}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return Result{Status: "UP", Latency: latency}
	}
	return Result{Status: "DOWN", Latency: latency}
}
