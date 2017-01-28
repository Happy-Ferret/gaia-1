package server

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//"github.com/notyim/gaia/types"
	"github.com/notyim/gaia/client"
	"github.com/notyim/gaia/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	for _, check := range h.server.Checks {
		fmt.Fprintf(w, fmt.Sprintf("%s,%s,%s\n", check.ID.Hex(), check.URI, check.Type))
	}
}

// Register a client to our internal server state
func (h *HTTPServer) ListClient(w http.ResponseWriter, r *http.Request) {
	lines := ""
	for _, c := range h.server.Clients {
		log.Printf("Existing clients %v\n", h.server.Clients)
		lines += fmt.Sprintf("IP: %s Location: %s\n", c.IpAddress, c.Location)
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
		IpAddress: ip,
		Location:  location,
	}
	log.Println("Found %s %s", ip, location)
	for _, c := range h.server.Clients {
		if c.IpAddress == client.IpAddress {
			w.WriteHeader(http.StatusAlreadyReported)
			return
		}
	}
	h.server.Clients = append(h.server.Clients, &client)

	log.Printf("Existing clients %v\n", h.server.Clients)
	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPServer) Install(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "#!/bin/bash\necho install run")
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPServer) Stats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Gaia Ok\n.Install server with.\ncurl https://gaia.noty.im/install | bash")
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPServer) CreateCheckResult(w http.ResponseWriter, r *http.Request) {
	isError := (r.FormValue("Error") == "true")
	errorMessage := r.FormValue("ErrorMessage")

	totalTime, err1 := strconv.Atoi(r.FormValue("TotalTime"))
	totalSize, err2 := strconv.Atoi(r.FormValue("TotalSize"))
	if err1 != nil || err2 != nil {
		log.Println("Fail to parse int", err1, err2)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "OK")
		return
	}

	checkedAt, err := strconv.Atoi(r.FormValue("CheckedAt"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Fail to parse CheckedAt", err)
		fmt.Fprintf(w, "OK")
		return
	}

	vars := mux.Vars(r)
	checkResult := models.CheckResult{
		CheckedAt:    time.Unix(0, int64(checkedAt)),
		CheckID:      vars["id"],
		Error:        isError,
		ErrorMessage: errorMessage,
		TotalTime:    time.Duration(int64(time.Millisecond) * int64(totalTime)),
		TotalSize:    totalSize,
	}
	//TODO: We should use a consumer channerl to process this instead of spawing goroutine instantly
	go checkResult.Save()

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
