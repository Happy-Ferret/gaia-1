package types

import (
	"time"
)

type Check struct {
	ID   string
	URI  string
	Type string
}

type HTTPCheckResponse struct {
	CheckedAt    time.Time
	CheckID      string
	TotalTime    time.Duration
	TotalSize    int
	HeaderSize   int
	BodySize     int
	Error        bool
	ErrorMessage string
}
