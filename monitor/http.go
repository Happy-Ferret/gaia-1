package monitor

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// HTTPServer represents internal http server
type HTTPServer struct {
	agent *Agent
	r     *mux.Router
}

// NewHTTPServer create a HttpServer struct
func NewHTTPServer(agent *Agent) *HTTPServer {
	s := &HTTPServer{
		agent: agent,
		r:     mux.NewRouter(),
	}
	s.r.HandleFunc("/", index)
	return s
}

// Start run http server
func (s *HTTPServer) Start() {
	log.Printf("Start server bootstrap")
	http.ListenAndServe("127.0.0.1:23501", s.r)
	log.Printf("Finish server bootstrap")
}

func index(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Gaia is running")
}
