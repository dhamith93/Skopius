package models

import "time"

type CheckResult struct {
	AgentID   string    `json:"agent_id"`
	Service   string    `json:"service"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	Code      int       `json:"code"`
	Latency   int64     `json:"latency"` // ms
	Error     string    `json:"error"`
	Timestamp time.Time `json:"time"`
	Received  time.Time `json:"received"`
}
