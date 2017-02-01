package types

import (
	"time"
)

type Check struct {
	ID       string
	URI      string
	Type     string
	Interval time.Duration
}

type HTTPCheckResponse struct {
	CheckedAt    time.Time                `json:"checked_at"`
	CheckID      string                   `json:"check_id"`
	Time         map[string]time.Duration `json:"time"`
	Body         string                   `json:"body"`
	Http         map[string]string        `json:"http"`
	Headers      map[string]string        `json:"headers"`
	Tcp          map[string]string        `json:"tcp"`
	Error        bool                     `json:"error"`
	ErrorMessage string                   `json:"error_message"`
}
