package monitor

import (
	"fmt"
	"github.com/notyim/gaia/config"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Agent constants.
const (
	AgentCapacity    = 600     // How many job agent can handle
	AgentBatchCheck  = 30      // How many check run in same go-routine
	AgentSignalStart = "start" // signal start string
	AgentSignalStop  = "stop"  // signal stop string
)

// Agent represent an agent that run checks
type Agent struct {
	Config     *config.Config
	InChan     chan Service
	out        chan StatusResult
	sigChan    chan string
	services   []Service
	httpClient *http.Client
}

// NewAgent create an agent with passing Output channel
func NewAgent(Out chan StatusResult) (*Agent, error) {
	a := &Agent{}
	a.InChan = make(chan Service, AgentCapacity)
	a.out = Out
	a.services = make([]Service, AgentCapacity, AgentCapacity)
	a.sigChan = make(chan string)

	tr := &http.Transport{
		//TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: false,
		DisableKeepAlives:  false,
	}

	a.httpClient = &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Transport: tr,
		Timeout:   time.Duration(10) * time.Second,
	}
	return a, nil
}

// Collect run check and record result
func (a *Agent) Collect() {
	log.Printf("Collect data")
	services := a.services

	var wg sync.WaitGroup

	for i := 0; i < AgentCapacity/AgentBatchCheck; i++ {
		j := i * AgentBatchCheck
		k := j + AgentBatchCheck + 1
		if k >= AgentCapacity {
			k = AgentCapacity
		}

		wg.Add(1)
		go func(batch []Service) {
			defer wg.Done()
			for _, s1 := range batch {
				if s1.Address != "" {
					log.Printf("Fetch for %s", s1.Address)
					a.out <- a.fetch(&s1)
				}
			}
		}(services[j:k])
	}
	wg.Wait()
}

// Start accepts data from input and sig channel for control flow
func (a *Agent) Start() {
	var total = 0
	for {
		select {
		case s := <-a.InChan:
			a.services[total] = s
			total++
		case s := <-a.sigChan:
			if s == AgentSignalStop {
				break
			}
		}
	}
}

// Stop signals agent to stop
func (a *Agent) Stop() {
	a.sigChan <- AgentSignalStop
}

func (a *Agent) fetch(s *Service) StatusResult {
	start := time.Now()
	rs := StatusResult{}
	rs.Service = s

	req, err := http.NewRequest("GET", s.Address, nil)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Printf("Error %v for %s", err, s.Address)
		rs.Response.Error = err
		rs.Response.Status = -1
	} else {
		rs.Response.Status = resp.StatusCode
		rs.Response.Duration = time.Since(start)
		body, _ := ioutil.ReadAll(resp.Body)
		rs.Response.Body = fmt.Sprintf("%s", body)
		resp.Body.Close()
	}
	return rs
}
