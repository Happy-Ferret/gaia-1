package monitor

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type HttpServer struct {
	agent *Agent
	r     *mux.Router
}

func NewHttpServer(agent *Agent) *HttpServer {
	s := &HttpServer{
		agent: agent,
		r:     mux.NewRouter(),
	}
	s.r.HandleFunc("/", index)
	return s
}

func (s *HttpServer) Start() {
	http.ListenAndServe("127.0.0.1:23501", s.r)
}

func index(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Gaia is running")
}
