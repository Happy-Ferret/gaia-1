package client

import (
	"github.com/notyim/gaia/types"
	//"io/ioutil"
	"bufio"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Only support 10k check atm
const MaxChecks = 10000

type Scanner struct {
	Checks         []*types.Check
	GaiaServerHost string
	m              *sync.Mutex
}

func NewScanner(gaiaServerHost string, ch chan string) *Scanner {
	s := Scanner{
		Checks:         make([]*types.Check, MaxChecks),
		GaiaServerHost: gaiaServerHost,
		m:              &sync.Mutex{},
	}

	go s.Listen(ch)
	go s.Monitor()
	return &s
}

func (s *Scanner) Listen(ch chan string) {
	for c := range ch {
		check := decode(c)
	}
}

// Connect to Gaia server to get initial check list
// TODO: We will switch to a TCP server, it makes thing much simpler
func (s *Scanner) Sync() {
	resp, err := http.Get(s.GaiaServerHost)
	if err != nil {
		log.Fatalf("Cannot sync initialize check %v", err)
	}
	defer resp.Body.Close()

	//body, err := ioutil.ReadAll(resp.Body)
	lineScanner := bufio.NewScanner(resp.Body)
	i := 0
	for lineScanner.Scan() {
		if line := lineScanner.Text(); line != "" {
			if check := decode(line); check != nil {
				s.Checks[i] = check
			}
		}
	}
}

func (s *Scanner) AddCheck(check *types.Check) {
}

// Monitor query stat of all checks in its internal data stucture and
// foward back to gaia server
func (s *Scanner) Monitor() {
}

func decode(s string) *types.Check {
	parts := strings.Split(s, ",")

	return &types.Check{
		ID:   parts[0],
		URI:  parts[1],
		Type: parts[2],
	}
}
