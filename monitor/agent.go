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
	AGENT_CAPACITY    = 600
	AGENT_BATCH_CHECK = 30
)

type Agent struct {
	Config     *config.Config
	InChan     chan Service
	out        chan StatusResult
	services   []Service
	httpClient *http.Client
}

func NewAgent(Out chan StatusResult) (*Agent, error) {
	a := &Agent{}
	a.InChan = make(chan Service, AGENT_CAPACITY)
	a.out = Out
	a.services = make([]Service, AGENT_CAPACITY, AGENT_CAPACITY)

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

func (a *Agent) Collect() {
	log.Printf("Collect data")

	var wg sync.WaitGroup

	for i := 0; i < AGENT_CAPACITY/AGENT_BATCH_CHECK; i++ {
		j := i * AGENT_BATCH_CHECK
		k := j + AGENT_BATCH_CHECK + 1
		if k >= AGENT_CAPACITY {
			k = AGENT_CAPACITY
		}

		wg.Add(1)
		log.Printf("Fetch batch %d", i)
		go func(batch []Service) {
			defer wg.Done()
			for _, s1 := range batch {
				if s1.Address != "" {
					log.Printf("Fecth for %s", s1.Address)
					a.out <- a.fetch(&s1)
				}
			}
		}(a.services[j:k])
	}
	wg.Wait()
}

func (a *Agent) Start() {
	var total = 0
	for {
		s := <-a.InChan
		a.services[total] = s
		total += 1
	}
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
