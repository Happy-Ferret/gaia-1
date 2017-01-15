package client

import (
	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	r *mux.Router
}

// Run the whole client
// Register route, initalize client, syncing
func CreateHTTPServer(scanner *Scanner) *HTTPServer {
	s := HTTPServer{
		r: mux.NewRouter(),
	}

	SyncCheck()
	//s.r.Handler("/checks").method("POST")
	return &s
}

// Retrieve existing checks to register to agen
func SyncCheck() {
}

// Handler
func RegisterHandler() {
}
