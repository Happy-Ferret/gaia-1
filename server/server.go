package server

import (
	"fmt"
	"github.com/notyim/gaia/client"
	"github.com/notyim/gaia/config"
	"github.com/notyim/gaia/db/mongo"
	"github.com/notyim/gaia/models"
	"log"
	"os"
	"os/signal"
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

	// Setup go routine for periodically sync

}

// Initialize gaia server
func Start(c *config.Config) {
	mongo.Connect("127.0.0.1:27017", "trinity_development")

	bindTo := fmt.Sprintf("%s:%d", "0.0.0.0", 28300)
	log.Println("Initalize server and bind to", bindTo)

	s := NewServer(c)
	s.SyncChecks()
	go s.HTTPServer.Start(bindTo)

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
