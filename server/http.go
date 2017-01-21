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
	fmt.Fprintf(w, "OK")
}

// Register a client to our internal server state
func (h *HTTPServer) ListClient(w http.ResponseWriter, r *http.Request) {
	for _, c := range h.server.Clients {
		log.Printf("Existing clients %v\n", h.server.Clients)
		fmt.Fprintf(w, fmt.Sprintf("IP: %s Location: %s\n", c.IpAddress, c.Location))
	}
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

func (h *HTTPServer) CreateCheckResult(w http.ResponseWriter, r *http.Request) {
	isError := (r.FormValue("Error") == "true")
	errorMessage := r.FormValue("ErrorMessage")

	totalTime, err1 := strconv.Atoi(r.FormValue("TotalTime"))
	totalSize, err2 := strconv.Atoi(r.FormValue("TotalSize"))
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "OK")
		return
	}

	checkResult := models.CheckResult{
		Error:        isError,
		ErrorMessage: errorMessage,
		TotalTime:    time.Duration(totalTime) * time.Second,
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

	s.r.HandleFunc("/checks", s.ListCheck).Methods("GET")
	s.r.HandleFunc("/client/register", s.RegisterClient).Methods("POST")
	s.r.HandleFunc("/clients", s.ListClient).Methods("GET")
	s.r.HandleFunc("/check_results", s.CreateCheckResult).Methods("POST")
	return &s
}
