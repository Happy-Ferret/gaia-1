package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/notyim/gaia/client"
	"github.com/notyim/gaia/models"
	"github.com/notyim/gaia/types"
	"log"
	"net/http"
	"os"
	//"strconv"
	//"time"
)

type HTTPServer struct {
	server  *Server
	r       *mux.Router
	flusher *Flusher
}

func (h *HTTPServer) Start(bindTo string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, h.r)
	if err := http.ListenAndServe(bindTo, loggedRouter); err != nil {
		log.Fatal("Error: Gaia HTTP server error", err)
	}
}

// Handle register a new check on http interface
func (h *HTTPServer) ListCheck(w http.ResponseWriter, r *http.Request) {
	for _, check := range h.server.Checks {
		fmt.Fprintf(w, fmt.Sprintf("%s,%s,%s\n", check.ID.Hex(), check.URI, check.Type))
	}
}

// Register a client to our internal server state
func (h *HTTPServer) ListClient(w http.ResponseWriter, r *http.Request) {
	lines := ""
	for _, c := range h.server.Clients {
		log.Printf("Existing clients %v\n", h.server.Clients)
		lines += fmt.Sprintf("IP: %s Location: %s\n", c.Address.IpAddress, c.Address.Location)
	}
	fmt.Fprintf(w, lines)
	w.WriteHeader(200)
}

// Register a client to our internal server state
func (h *HTTPServer) RegisterClient(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")
	location := r.FormValue("location")

	if ip == "" || location == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := client.Client{
		Address: types.ClientAddress{
			IpAddress: ip,
			Location:  location,
		},
	}
	for _, c := range h.server.Clients {
		if c.Address.IpAddress == client.Address.IpAddress {
			w.WriteHeader(http.StatusAlreadyReported)
			return
		}
	}
	h.server.Clients = append(h.server.Clients, &client)

	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPServer) Install(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("#!/bin/bash\necho Start Installer\ncurl -s https://github.com/NotyIm/gaia/releases/download/%s/gaia_am64.deb -o /tmp/gaia.deb\nsudo dpkg -i /tmp/gaia.deb\nsystemctl enable gaia\nsystemctl start gaia\nrm /tmp/gaia.deb", types.Version))
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPServer) Stats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Gaia Ok\n.Install server with.\ncurl -s https://gaia.noty.im/install | bash")
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPServer) CreateCheckResult(w http.ResponseWriter, r *http.Request) {
	checkResult := models.CheckResult{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&checkResult)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Fail to parse Check Result body", err)
		fmt.Fprintf(w, "Fail to parse Check Result body", err)
		return
	}

	h.flusher.Write(&checkResult)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "OK")
}

// Run the whole client
// Register route, initalize client, syncing
func CreateHTTPServer(server *Server, flusher *Flusher) *HTTPServer {
	s := HTTPServer{
		server:  server,
		flusher: flusher,
		r:       mux.NewRouter(),
	}

	s.r.HandleFunc("/", s.Stats).Methods("GET")
	s.r.HandleFunc("/install", s.Install).Methods("GET")
	s.r.HandleFunc("/checks", s.ListCheck).Methods("GET")
	s.r.HandleFunc("/client/register", s.RegisterClient).Methods("POST")
	s.r.HandleFunc("/clients", s.ListClient).Methods("GET")
	s.r.HandleFunc("/check_results/{id}", s.CreateCheckResult).Methods("POST")
	return &s
}
