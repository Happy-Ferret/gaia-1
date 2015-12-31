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
	Service  *Service
	Response struct {
		Body     string
		Status   int
		Duration time.Duration
		Error    error
	}
}

func NewService(Address, Id string) Service {
	s := Service{Address, Id}
	return s
}
