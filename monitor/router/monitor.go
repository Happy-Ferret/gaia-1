package router

import (
	"encoding/json"
	"fmt"
	"github.com/notyim/gaia/monitor/core"
	"net/http"
)

func SaveMonitor(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t core.Service

	err := decoder.Decode(&t)
	if err != nil {
		// sth here
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("%v", t)
}

func UpdateMonitor(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Gaia is running")
}

func DeleteMonitor(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Gaia is running")
}

func GetService() {
}
