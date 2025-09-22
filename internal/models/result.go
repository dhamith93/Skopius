package models

import "time"

type CheckResult struct {
	ID      int64  `json:"id" db:"id"`
	AgentID string `json:"agent_id" db:"agent_id"`
	Service string `json:"service" db:"service"`
	URL     string `json:"url" db:"url"`
	Status  string `json:"status" db:"status"`
	Code    int    `json:"code" db:"code"`
	Error   string `json:"error" db:"error"`

	// Normalized durations (milliseconds)
	DNS     int64 `json:"dns" db:"dns"`
	Connect int64 `json:"connect" db:"connect"`
	TLS     int64 `json:"tls" db:"tls"`
	TTFB    int64 `json:"ttfb" db:"ttfb"` // Time To First Byte
	Server  int64 `json:"server" db:"server"`
	Total   int64 `json:"total" db:"total"`

	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Received  time.Time `json:"received" db:"received"`
}
