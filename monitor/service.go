package monitor

import (
	"time"
)

// Service is a logic unit that needs to be monitored.
// It can be:
// an URL, an Ip address
type Service struct {
	Address string // address of service
	ID      string // id of service
}

// StatusResult represent check data of a service
type StatusResult struct {
	Service  *Service
	Response struct {
		Body     string
		Status   int
		Duration time.Duration
		Error    error
	}
}

// NewService creates a service struct
func NewService(Address, ID string) Service {
	s := Service{Address, ID}
	return s
}
