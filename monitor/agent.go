package monitor

import (
	"github.com/notyim/gaia/config"
)

type Agent struct {
	Config *config.Config
}

func NewAgent() (*Agent, error) {
	a := &Agent{}
	return a, nil
}

// Main entry point for monitoring system
func Start() {

}
