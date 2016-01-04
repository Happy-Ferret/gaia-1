package monitor

import (
	"log"
	"os"
	"os/signal"
	//"time"
	"github.com/yvasiyarov/gorelic"
)

// Monitor constants
const (
	CheckInterval = 1000 // time interval between checks
)

// Start is the main entry point for monitoring system
func Start() {
	log.Printf("Initalize...")

	registerNewRelic()

	shutdown := make(chan struct{})

	flusher := NewFlusher(nil)
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

	log.Printf("Register handle")
	registerSignal(shutdown)

	log.Printf("Start monitoring")
	registerMonitor(agent, shutdown)

	go registerHTTPServer(agent)
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
	status := make(chan string)

	for {
		// Run collect in go-routine to give back execution control to main thread
		// so when we ctrl-c, it can be catched instantly and exit
		go func(ch chan string) {
			agent.Collect()
			log.Printf("Done collect")
			ch <- "done"
		}(status)

		select {
		case s := <-shutdown:
			log.Printf("Receive shutdown command %v. Will quit", s)
			//@TODO cleanup when exiting
			os.Exit(1)
			return
		case s := <-status:
			log.Printf("Get status %v", s)
			if s == "done" {
				continue
			} else {
				log.Fatalf("Bad signal %s", s)
			}
			//case t := <-ticker.C:
			//	log.Printf("Tick at %v", t)
			//	continue
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
