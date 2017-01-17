package server

import (
	"fmt"
	"github.com/notyim/gaia/client"
	"github.com/notyim/gaia/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	bindTo := "0.0.0.0:28300"
	log.Println("Initalize server and bind to")

	flusher := NewFlush()
	s := NewServer(c, f)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	userSignal := <-sigChan
	log.Println("Got signal", userSignal, "from end-user")
	log.Println("Attempt to quit")
	os.Exit(0)
}

func NewServer(c *config.Config, f *Flusher) *Server {
	s := Server{
		Clients: make([]*client.Client, InitClientsSize),
		config:  c,
	}

	h := CreateHTTPServer(f)
	s.HTTPServer = h

	go h.Start()

	return &s
}
