package server

import (
	"fmt"
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
}

// Sync and keep track of checks from db
// This is poorman changedfeed in MongoDB
// I wish we can use RethinkDB here
func (s *Server) SyncChecks() {
	// Initial sync
	//s.Checks.All()
	var checks []models.Check
	models.AllChecks(&checks)
	s.Checks = checks

	log.Println("Found check %v", s.Checks)

	ticker := time.NewTicker(time.Second * 3)
	// Setup go routine for periodically sync
	go func() {
		for t := range ticker.C {
			log.Println("Poll checks at", t)
			var checks []models.Check
			models.FindChecksAfter(&checks, s.Checks[len(s.Checks)-1].ID)
			s.Checks = append(s.Checks, checks...)
			for _, check := range checks {
				s.PushCheckToClients(&check)
			}
		}
	}()
}

//Push new checks to client
func (s *Server) PushCheckToClients(check *models.Check) {
	for _, c := range s.Clients {
		// Implement https for client
		// TODO We will dismiss all this and replica with a TCP with tls
		_, err := http.PostForm(fmt.Sprintf("http://%s:28302/checks", c.IpAddress),
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
		config: c,
	}

	h := CreateHTTPServer(&s, NewFlusher())
	s.HTTPServer = h

	return &s
}
