package router

import (
	"encoding/json"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
	"github.com/influxdb/influxdb/client/v2"
	"github.com/notyim/gaia/env"
	"github.com/notyim/gaia/monitor/core"
	"github.com/notyim/gaia/types"
	"log"
	"net/http"
)

//@TODO Refactor
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	env := env.GetEnv()
	q := client.Query{
		Command:  cmd,
		Database: env.Config.InfluxdbDb,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

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

func GetService(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serviceId := vars["id"]
	log.Printf("Service Request %d", serviceId)

	/*
		client := env.GetEnv().Influx
		q := fmt.Sprintf("select * from \"24h\".http_response where ServiceId='%d' order by time desc limit 1 ", service)
		res, err := queryDB(client, q)
		if err != nil {
			resp.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(resp, "Cannot fetch service")
			return
		}
		for i, row := range res[0].Series[0].Values {
			t, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				log.Println(err)
				resp.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(resp, "Server error")
				return
			}
			log.Printf("Row = %v", row)
			log.Printf("Get Service [%2d] %s: %s\n", i, t.Format(time.Stamp), row)
			b, err := json.Marshal(row)
			fmt.Fprintf(resp, "%s", b)
		}
	*/

	session := env.GetEnv().Rethink

	res, err := r.Table("http_response").Get(serviceId).Run(session)
	if err != nil {
		//@TODO return JSON
		log.Printf("Error when fetching http_response %v", err)
		fmt.Fprintf(resp, "Error")
		return
	}
	defer res.Close()

	var doc types.RethinkService
	err = res.One(&doc)
	if err == r.ErrEmptyResult {
		fmt.Fprintf(resp, "not found")
		return
	}
	b, err := json.Marshal(doc)
	log.Println(b)
	log.Println(err)
	fmt.Fprintf(resp, string(b))
}
