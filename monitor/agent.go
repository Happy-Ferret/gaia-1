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

const (
	AGENT_CAPACITY = 1000
)

type Agent struct {
	Config   *config.Config
	InChan   chan Service
	out      chan StatusResult
	services []Service
}

func NewAgent(Out chan StatusResult) (*Agent, error) {
	a := &Agent{}
	a.InChan = make(chan Service, AGENT_CAPACITY)
	a.out = Out
	a.services = make([]Service, 40, AGENT_CAPACITY)
	return a, nil
}

func (a *Agent) Collect() {
	log.Printf("Collect data")

	var wg sync.WaitGroup

	for _, s := range a.services {
		if s.Address != "" {
			wg.Add(1)
			go func(s1 Service) {
				log.Printf("Check service add %v", s1.Address)
				log.Printf("Check service %v", s1)
				defer wg.Done()
				a.out <- a.fetch(&s1)
			}(s)
		}
	}
	wg.Wait()
}

func (a *Agent) Start() {
	var total = 0
	for {
		s := <-a.InChan
		a.services[total] = s
		log.Printf("Got service", s)
		log.Printf("Total", total)
		total += 1
	}
}

func (a *Agent) fetch(s *Service) StatusResult {
	log.Printf("Will fetch %v", s)
	start := time.Now()
	rs := StatusResult{}

	tr := &http.Transport{
		//TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: false,
	}

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
		Transport: tr,
	}
	req, err := http.NewRequest("GET", s.Address, nil)

	resp, err := client.Do(req)
	if err != nil {
		rs.Error = err
	} else {
		rs.Response.Status = resp.StatusCode
		rs.Response.Duration = time.Since(start)
		body, _ := ioutil.ReadAll(resp.Body)
		rs.Response.Body = fmt.Sprintf("%s", body)
		resp.Body.Close()
	}
	return rs
}
