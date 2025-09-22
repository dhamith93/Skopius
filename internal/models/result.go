package models

import "time"

type CheckResult struct {
	ID        int64     `json:"id" db:"id"`
	AgentID   string    `json:"agent_id" db:"agent_id"`
	Service   string    `json:"service" db:"service"`
	URL       string    `json:"url" db:"url"`
	Status    string    `json:"status" db:"status"`
	Code      int       `json:"code" db:"code"`
	Latency   int64     `json:"latency" db:"latency"`
	Error     string    `json:"error" db:"error"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Received  time.Time `json:"received" db:"received"`
}
