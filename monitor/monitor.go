package monitor

import (
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	CHECK_INTERVAL = 1000
)

// Main entry point for monitoring system
func Start() {
	log.Printf("Initalize...")
	shutdown := make(chan struct{})

	flusher := NewFlusher()
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
	ticker := time.NewTicker(CHECK_INTERVAL * time.Millisecond)

	for {
		// Run collect in go-routine to give back execution control to main thread
		// so when we ctrl-c, it can be catched instantly and exit
		go func() {
			agent.Collect()
		}()

		select {
		case s := <-shutdown:
			log.Printf("Receive shutdown command %v. Will quit", s)
			//@TODO cleanup when exiting
			os.Exit(1)
			return
		case t := <-ticker.C:
			log.Printf("Tick at %v", t)
			continue
		}
	}
}

func registerHttp() {
}
