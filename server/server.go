package server

import (
	"bytes"
	"fmt"
	"github.com/newrelic/go-agent"
	"github.com/notyim/gaia/apm"
	"github.com/notyim/gaia/client"
	"github.com/notyim/gaia/config"
	"github.com/notyim/gaia/db/influxdb"
	"github.com/notyim/gaia/db/mongo"
	"github.com/notyim/gaia/models"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
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
	httpClient *http.Client
}

// Sync and keep track of checks from db
// This is poorman changedfeed in MongoDB
// I wish we can use RethinkDB here
func (s *Server) SyncChecks() {
	// Setup go routine for periodically sync
	ticker := time.NewTicker(time.Second * 15)
	go func() {
		shard := 0

		for t := range ticker.C {
			shard += 1

			log.Println("Syncing shard", shard, "at", t)

			var checks []models.Check
			models.FindChecksByShard(&checks, shard)
			if checks != nil && len(checks) > 0 {
				s.PushBulkCheckToClients(checks)
			} else {
				log.Println("Found no check for shard", shard)
			}
			if shard >= 4 {
				shard = 0
			}
		}
	}()
}

func (s *Server) PushBulkCheckToClients(checks []models.Check) {
	txn := apm.NewrelicApp.StartTransaction("PushBulkCheckToClient", nil, nil)
	defer txn.End()

	lines := make([]string, len(checks))
	for i, check := range checks {
		lines[i] = fmt.Sprintf("%s,%s,%s", check.ID.Hex(), check.URI, check.Type)
	}
	payload := strings.Join(lines, "\n")

	for _, c := range s.Clients {
		defer newrelic.StartSegment(txn, "HttpClientPost").End()
		// TODO We will dismiss all this and replica with a TCP with tls
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:28302/bulkchecks", c.Address.IpAddress), bytes.NewBufferString(payload))
		if err != nil {
			log.Println("Error Fail to create http request", err)
			continue
		}
		_, err = s.httpClient.Do(req)
		if err != nil {
			log.Println("Error fail to push bulk checks to client", err)
		}
	}
}

//Push new checks to client
func (s *Server) PushCheckToClients(check *models.Check) {
	log.Println("Sync Check", check, "to client")
	for _, c := range s.Clients {
		txn := apm.NewrelicApp.StartTransaction("PushCheckToClient", nil, nil)
		defer txn.End()

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
