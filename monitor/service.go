package monitor

import (
	"time"
)

// A service is a logic unit that needs to be monitored.
// It can be:
// an URL, an Ip address
type Service struct {
	Address string
	Id      string
}

type StatusResult struct {
	Latency  int
	Error    error
	Response struct {
		Status   int
		Duration time.Duration
	}
}

func NewService(Address, Id string) Service {
	s := Service{Address, Id}
	return s
}
