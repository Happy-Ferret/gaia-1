package monitor

import (
	"github.com/notyim/gaia/config"
	"github.com/yvasiyarov/gorelic"
	"log"
	"os"
	"os/signal"
	//"time"
)

// Monitor constants
const (
	CheckInterval = 1000 // time interval between checks
)

// Start is the main entry point for monitoring system
func Start() {
	config := config.NewConfig()

	log.Printf("Initalize...")

	//registerNewRelic()

	shutdown := make(chan struct{})

	flusher := NewFlusher(config, nil)
	go func() {
		flusher.Start()
	}()

	agent, _ := NewAgent(flusher.DataChan)
	go func() {
		agent.Start()
	}()

	coordinator := NewCoordinator(agent.InChan)
	go func() {
		coordinator.Start()
	}()

	log.Printf("Register http server")
	go registerHTTPServer(agent)

	log.Printf("Register monitoring point")
	registerMonitor(agent, shutdown)

	log.Printf("Register signal handle")
	registerSignal(shutdown)
}

func registerSignal(shutdown chan struct{}) {
	log.Printf("Register signal")
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		s := <-signals
		log.Printf("Got signal %v", s)
		close(shutdown)
	}()
}

func registerMonitor(agent *Agent, shutdown chan struct{}) {
	// Every 3 seconds
	//ticker := time.NewTicker(CheckInterval * time.Millisecond)
	for {
		select {
		case s := <-shutdown:
			log.Printf("Receive shutdown command %v. Will quit", s)
			//@TODO cleanup when exiting
			os.Exit(1)
			return
		}
	}
}

func registerHTTPServer(agent *Agent) {
	h := NewHTTPServer(agent)
	h.Start()
}

func registerNewRelic() {
	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicLicense = os.Getenv("NEWRELIC_LICENSE")
	log.Printf("NRL %s", agent.NewrelicLicense)
	agent.NewrelicName = "Gaia"
	agent.CollectHTTPStat = true
	agent.Verbose = true

	agent.Run()
}
