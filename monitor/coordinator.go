package monitor

type Coordinator struct {
	AgentChan chan Service
}

func NewCoordinator(agentChan chan Service) *Coordinator {
	c := &Coordinator{agentChan}
	return c
}

func (c *Coordinator) Start() {
	// @TODO
	// Fetch data from source in a loop and notify agent channel about new data
	// or notify agent channel about removing of data
	c.AgentChan <- NewService("www.axcoto.com", "1")
}
