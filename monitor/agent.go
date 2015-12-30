package monitor

import (
	"github.com/notyim/gaia/config"
	"log"
)

type Agent struct {
	Config *config.Config
	InChan chan Service
	out    chan StatusResult
}

const (
	AGENT_CAPACITY = 1000
)

func NewAgent(Out chan StatusResult) (*Agent, error) {
	a := &Agent{}
	a.InChan = make(chan Service, AGENT_CAPACITY)
	a.out = Out
	return a, nil
}

func (a *Agent) Collect() {
	log.Printf("Collect data")
	a.out <- StatusResult{100, 100}
}

func (a *Agent) Start() {
	// First run, collect instantly any data
	a.Collect()

}
