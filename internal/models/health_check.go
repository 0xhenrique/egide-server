package models

import "time"

type HealthCheck struct {
	ID             int64     `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	ResponseTimeMs int       `json:"response_time_ms"`
	StatusCode     *int      `json:"status_code,omitempty"` // nil if request failed completely
	Success        bool      `json:"success"`
	Error          *string   `json:"error,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
