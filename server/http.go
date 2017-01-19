package server

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//"github.com/notyim/gaia/types"
	"github.com/notyim/gaia/client"
	"log"
	"net/http"
	"os"
)

type HTTPServer struct {
	server  *Server
	r       *mux.Router
	flusher *Flusher
}

func (h *HTTPServer) Start(bindTo string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, h.r)
	log.Println(http.ListenAndServe(bindTo, loggedRouter))
}

// Handle register a new check on http interface
func (h *HTTPServer) ListCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// Register a client to our internal server state
func (h *HTTPServer) RegisterClient(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")
	location := r.FormValue("location")

	client := client.Client{
		IpAddress: ip,
		Location:  location,
	}
	log.Println("Found %s %s", ip, location)
	h.server.Clients = append(h.server.Clients, &client)
	log.Printf("Existing clients %v\n", h.server.Clients)
}

// Run the whole client
// Register route, initalize client, syncing
func CreateHTTPServer(server *Server, flusher *Flusher) *HTTPServer {
	s := HTTPServer{
		server:  server,
		flusher: flusher,
		r:       mux.NewRouter(),
	}

	s.r.HandleFunc("/checks", s.ListCheck).Methods("GET")
	s.r.HandleFunc("/client/register", s.RegisterClient).Methods("POST")
	return &s
}
