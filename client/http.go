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
	scanner *Scanner
}

// Handle register a new check on http interface
func (h *HTTPServer) RegisterCheck(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	uri := r.FormValue("uri")
	checkType := r.FormValue("type")
	check := types.Check{id, uri, checkType}
	h.scanner.AddCheck(&check)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "OK")
}

// Run the whole client
// Register route, initalize client, syncing
func CreateHTTPServer(scanner *Scanner) *HTTPServer {
	s := HTTPServer{
		scanner: scanner,
		r:       mux.NewRouter(),
	}

	s.r.HandleFunc("/checks", s.RegisterCheck).Methods("POST")
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.r)
	log.Println(http.ListenAndServe("0.0.0.0:28302", loggedRouter))
	return &s
}
