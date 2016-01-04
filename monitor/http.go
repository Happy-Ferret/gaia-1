package monitor

import (
	"fmt"
	"github.com/gorilla/mux"
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
	http.ListenAndServe("127.0.0.1:23501", s.r)
}

func index(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Gaia is running")
}
