package monitor

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Service struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
}

type CheckResult struct {
	Name      string
	URL       string
	Status    string // "UP" or "DOWN"
	Code      int
	Latency   int64 // ms
	Timestamp time.Time
	Error     string
}

func (s *Service) Check() CheckResult {
	client := http.Client{Timeout: 5 * time.Second}
	start := time.Now()
	resp, err := client.Get(s.URL)
	latency := time.Since(start).Milliseconds()

	result := CheckResult{
		Name:      s.Name,
		URL:       s.URL,
		Latency:   latency,
		Timestamp: time.Now(),
	}

	if err != nil {
		result.Status = "DOWN"
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.Code = resp.StatusCode
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Status = "UP"
	} else {
		result.Status = "DOWN"
	}
	return result
}

func (s *Service) Probe(ctx context.Context, results chan<- CheckResult) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			res := s.Check()
			results <- res
		case <-ctx.Done():
			log.Printf("Probe stopped for %s\n", s.Name)
			return
		}
	}
}
