package monitor

// A service is a logic unit that needs to be monitored.
// It can be:
// an URL, an Ip address
type Service struct {
	Address string
	Id      string
}

type StatusResult struct {
	Latency int
	Status  int
}

func NewService(Address, Id string) Service {
	s := Service{Address, Id}
	return s
}
