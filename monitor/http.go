package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
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

	stat := stats.New()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		s, err := json.Marshal(stat.Data())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(s)
	})
	s.r.HandleFunc("/_stats", h)

	return s
}

// Start run http server
func (s *HTTPServer) Start() {
	log.Printf("Start server bootstrap")
	bind := "127.0.0.1:23501"
	if cbind := os.Getenv("BIND"); cbind != "" {
		bind = cbind
	}
	log.Println(http.ListenAndServe(bind, s.r))
	log.Printf("Finish server bootstrap")
}

func index(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Gaia is running")
}
