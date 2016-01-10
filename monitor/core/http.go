package core

import (
	"time"
)

// Service is a logic unit that needs to be monitored.
// It can be:
// an URL, an Ip address
type Service struct {
	Address  string // address of service
	ID       string // id of service
	Interval int64  // interval to run check for this service
	Type     string // type of service: http, tcp, ttl...
}

type ResponseMetric struct {
	Body     string
	Status   int
	Duration time.Duration
	Error    error
}

// HTTPMetric represent check data of a service
type HTTPMetric struct {
	Service  *Service
	Response ResponseMetric
}

// NewHTTPService creates a service struct
func NewHTTPService(Address, ID string, interval int64) *Service {
	s := &Service{Address, ID, interval, "http"}
	return s
}
