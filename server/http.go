package client

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/notyim/gaia/types"
	"log"
	"net/http"
	"os"
)

type HTTPServer struct {
	r       *mux.Router
	flusher *Flusher
}

func (h *HTTPServer) Start() {
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.r)
	log.Println(http.ListenAndServe("0.0.0.0:28300", loggedRouter))
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
		ip:       ip,
		location: loction,
	}
	h.s.append(client)
	h.s.Clients = append(h.s.Clients, client)
}

// Run the whole client
// Register route, initalize client, syncing
func CreateHTTPServer(flusher *Flusher) *HTTPServer {
	s := HTTPServer{
		flusher: flusher,
		r:       mux.NewRouter(),
	}

	s.r.HandleFunc("/checks", s.ListCheck).Methods("GET")
	s.r.HandleFunc("/client/register", s.RegisterClient).Methods("POST")
	return &s
}
