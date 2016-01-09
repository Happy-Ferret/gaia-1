package monitor

import (
	"fmt"
	"github.com/notyim/gaia/config"
	"github.com/notyim/gaia/monitor/core"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Agent constants.
const (
	AgentCapacity      = 600       // How many job agent can handle
	AgentBatchCheck    = 2         // How many check run in same go-routine
	AgentSignalStart   = "start"   // signal start string
	AgentSignalStop    = "stop"    // signal stop string
	AgentSignalCollect = "collect" // signal collect
)

// Agent represent an agent that run checks
type Agent struct {
	Config     *config.Config
	InChan     chan *core.Service
	out        chan *core.HTTPMetric
	sigChan    chan string
	services   []*core.Service
	pools      []chan string
	httpClient *http.Client
	lock       *sync.Mutex
}

// NewAgent create an agent with passing Output channel
func NewAgent(Out chan *core.HTTPMetric) (*Agent, error) {
	a := &Agent{
		lock: &sync.Mutex{},
	}
	a.InChan = make(chan *core.Service, AgentCapacity)
	a.out = Out
	a.services = make([]*core.Service, AgentCapacity, AgentCapacity)
	a.pools = make([]chan string, AgentCapacity, AgentCapacity)
	a.sigChan = make(chan string)

	tr := &http.Transport{
		//TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: false,
		DisableKeepAlives:  false,
	}

	a.httpClient = &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Transport: tr,
		Timeout:   time.Duration(15) * time.Second,
	}
	return a, nil
}

// Start accepts data from input and sig channel for control flow
func (a *Agent) Start() {
	var total = 0
	for {
		select {
		case s := <-a.InChan:
			go func() {
				a.lock.Lock()
				ch := make(chan string)
				a.services[total] = s
				a.pools[total] = ch
				total++
				a.lock.Unlock()
				a.newWorker(s, ch)
			}()
		case s := <-a.sigChan:
			if s == AgentSignalStop {
				// destroy pool
				a.destroyWorkers()
				break
			}
		}
	}
}

func (a *Agent) destroyWorkers() {
	for i, ch := range a.pools {
		if ch != nil {
			log.Printf("Stop worker %i\n", i)
			ch <- "stop"
		}
	}
}

func (a *Agent) newWorker(s *core.Service, ch chan string) error {
	timer := time.NewTicker(time.Duration(s.Interval) * time.Millisecond)

	select {
	case t := <-timer.C:
		log.Printf("Fetch for %s at %v", s.Address, t)
		// @TODO error handle and logging with Raven maybe?
		a.out <- a.fetch(s)
	case action := <-ch:
		// @TODO more action here
		log.Println("Got signal %s Quit worker service %s", action, s.Address)
		break
	}
	return nil
}

// Stop signals agent to stop
func (a *Agent) Stop() {
	a.sigChan <- AgentSignalStop
}

func (a *Agent) fetch(s *core.Service) *core.HTTPMetric {
	start := time.Now()
	rs := &core.HTTPMetric{}
	rs.Service = s

	req, err := http.NewRequest("GET", s.Address, nil)
	// Make sure we close http connection to avoid leaking file descriptor
	req.Close = true

	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Printf("Error %v for %s", err, s.Address)
		rs.Response.Error = err
		rs.Response.Status = -1
	} else {
		rs.Response.Error = err
		rs.Response.Status = resp.StatusCode
		rs.Response.Duration = time.Since(start)
		body, _ := ioutil.ReadAll(resp.Body)
		rs.Response.Body = fmt.Sprintf("%s", body)
		resp.Body.Close()
	}
	log.Printf("%s: %v", s.Address, rs.Response)
	return rs
}
