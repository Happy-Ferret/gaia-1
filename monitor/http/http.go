package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/notyim/gaia/monitor"
	"net/http"
)

type HttpServer struct {
	agent *monitor.Agent
	mux   *http.ServeMux
}

func NewHttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
}

func index(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Gaia is running")
}
