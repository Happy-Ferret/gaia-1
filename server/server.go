package server

import (
	"fmt"
	"github.com/notyim/gaia/client"
	"github.com/notyim/gaia/config"
	"github.com/notyim/gaia/db/influxdb"
	"github.com/notyim/gaia/db/mongo"
	"github.com/notyim/gaia/models"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

const (
	InitClientsSize = 10
)

type Server struct {
	Clients    []*client.Client
	Checks     []models.Check
	config     *config.Config
	HTTPServer *HTTPServer
	httpClient *net.HTTP
}

// Sync and keep track of checks from db
// This is poorman changedfeed in MongoDB
// I wish we can use RethinkDB here
func (s *Server) SyncChecks() {
	var checks []models.Check
	models.AllChecks(&checks)
	s.Checks = checks

	ticker := time.NewTicker(time.Second * 15)
	// Setup go routine for periodically sync
	go func() {
		shard := 0

		for t := range ticker.C {
			shard += 1

			log.Println("Syncing at", t, "for shard", shard)

			var checks []models.Check
			models.FindChecksByShard(&checks, shard)
			if checks != nil && len(checks) > 0 {
				s.PushBulkCheckToClients(checks)
			}
			if shard >= 4 {
				shard = 0
			}
		}
	}()
}

func (s *Server) PushBulkCheckToClients(checks []models.Check) {
	lines := make([]string, len(checks))
	for i, check := range checks {
		lines[i] = fmt.Sprintf("%s,%s,%s", check.ID, check.URI, check.Type)
	}
	payload := strings.Join(lines, "\n")
	for _, c := range s.Clients {
		// TODO We will dismiss all this and replica with a TCP with tls
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:28302/bulkchecks", c.Address.IpAddress), bytes.NewBufferString(payload))
		if err != nil {
			log.Println("Error Fail to create http request", err)
			continue
		}
		_, err := s.httpClient.Do(req)
		if err != nil {
			log.Println("Error fail to push bulk checks to client", err)
		}
	}

}

//Push new checks to client
func (s *Server) PushCheckToClients(check *models.Check) {
	log.Println("Sync Check", check, "to client")
	for _, c := range s.Clients {
		// Implement https for client
		// TODO We will dismiss all this and replica with a TCP with tls
		_, err := http.PostForm(fmt.Sprintf("http://%s:28302/checks", c.Address.IpAddress),
			url.Values{"id": {check.ID.Hex()}, "uri": {check.URI}, "type": {check.Type}})
		log.Println("Push", check, "to client", c)
		if err != nil {
			log.Println("Error Fail to push check to client", err)
		}
	}
}

// Initialize gaia server
func Start(c *config.Config) {
	mongo.Connect("127.0.0.1:27017", c.MongoDBName)

	log.Println("Initalize server and bind to", c.GaiaServerBindTo)

	influxdb.Connect(c.InfluxdbHost, c.InfluxdbUsername, c.InfluxdbPassword)
	influxdb.UseDB(c.InfluxdbDb)

	s := NewServer(c)
	s.SyncChecks()
	go s.HTTPServer.Start(c.GaiaServerBindTo)

	//@TODO Move this to config
	CreateWorker("localhost:6379", "0", "30")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	userSignal := <-sigChan
	log.Println("Got signal", userSignal, "from end-user")
	log.Println("Attempt to quit")
	os.Exit(0)
}

func NewServer(c *config.Config) *Server {
	s := Server{
		config:     c,
		httpClient: &http.Client{},
	}

	h := CreateHTTPServer(&s, NewFlusher())
	s.HTTPServer = h

	return &s
}
