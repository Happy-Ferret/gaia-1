package client

import (
	"bufio"
	"fmt"
	"github.com/notyim/gaia/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Only support 10k check atm
const MaxChecks = 10000

type Scanner struct {
	Checks             []*types.Check
	CheckInterval      time.Duration
	totalCheck         int
	GaiaServerHost     string
	m                  *sync.Mutex
	CheckRegisterQueue chan string
	f                  *Flusher
}

func NewScanner(gaiaServerHost string) *Scanner {
	s := Scanner{
		Checks:             make([]*types.Check, MaxChecks),
		totalCheck:         0,
		CheckInterval:      15 * time.Second,
		GaiaServerHost:     gaiaServerHost,
		m:                  &sync.Mutex{},
		CheckRegisterQueue: make(chan string),
		f:                  NewFlusher(gaiaServerHost),
	}

	s.Sync()
	go s.Listen()
	go s.Monitor()

	return &s
}

func (s *Scanner) Start() {
	f := NewFlusher(s.GaiaServerHost)
	f.Start()
}

func (s *Scanner) Listen() {
	for c := range s.CheckRegisterQueue {
		check := decode(c)
		s.Checks[s.totalCheck] = check
		s.totalCheck++
	}
}

// Connect to Gaia server to get initial check list
// TODO: We will switch to a TCP server, it makes thing much simpler
func (s *Scanner) Sync() {
	resp, err := http.Get(fmt.Sprintf("%s/checks", s.GaiaServerHost))

	if err != nil {
		log.Fatalf("Cannot sync initialize check %v", err)
	}
	defer resp.Body.Close()
	log.Println("Got initial set of checks. Start decoding checks result set.")

	lineScanner := bufio.NewScanner(resp.Body)
	for lineScanner.Scan() {
		if line := lineScanner.Text(); line != "" {
			if check := decode(line); check != nil {
				log.Println("Found sync check", check)
				s.Checks[s.totalCheck] = check
				s.totalCheck++
			}
		}
	}
}

func (s *Scanner) AddCheck(check *types.Check) {
	s.Checks[s.totalCheck] = check
	s.totalCheck++
}

func (s *Scanner) RemoveCheck() {
}

// Monitor query stat of all checks in its internal data stucture and
// foward back to gaia server
func (s *Scanner) Monitor() {
	ticker := time.Tick(s.CheckInterval)
	for tick := range ticker {
		log.Println("Execute check at", tick)
		for _, c := range s.Checks {
			if c == nil {
				continue
			}

			go s.Execute(c)
		}
	}
}

func (s *Scanner) Execute(check *types.Check) {
	log.Println("Evaluate", check.URI)

	startAt := time.Now()
	resp, err := http.Get(check.URI)
	endAt := time.Now()

	response := types.HTTPCheckResponse{
		CheckID:   check.ID,
		TotalTime: endAt.Sub(startAt),
	}
	if err != nil {
		log.Println(check.URI, "error", err)

		response.Error = false
		response.ErrorMessage = err.Error()
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			response.TotalSize = len(body)
			response.BodySize = len(body)
		}
	}
	s.f.Write(&response)
}

func decode(s string) *types.Check {
	parts := strings.Split(s, ",")

	return &types.Check{
		ID:   parts[0],
		URI:  parts[1],
		Type: parts[2],
	}
}
