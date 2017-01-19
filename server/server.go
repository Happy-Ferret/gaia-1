package server

import (
	"fmt"
	"github.com/notyim/gaia/client"
	"github.com/notyim/gaia/config"
	"log"
	"os"
	"os/signal"
)

const (
	InitClientsSize = 10
)

type Server struct {
	Clients    []*client.Client
	config     *config.Config
	HTTPServer *HTTPServer
}

// Initialize gaia server
func Start(c *config.Config) {
	bindTo := fmt.Sprintf("%s:%d", "0.0.0.0", 28300)
	log.Println("Initalize server and bind to", bindTo)

	s := NewServer(c)
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
