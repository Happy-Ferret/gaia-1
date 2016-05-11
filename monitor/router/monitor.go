package router

import (
	"encoding/json"
	"fmt"
	"github.com/influxdb/influxdb/client/v2"
	"github.com/notyim/gaia/env"
	"github.com/notyim/gaia/monitor/core"
	"log"
	"net/http"
	"time"
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
	/*
		c, _ := client.NewHTTPClient(client.HTTPConfig{
			Addr:     config.InfluxdbHost,
			Username: config.InfluxdbUsername,
			Password: config.InfluxdbPassword,
		})

		f.client = c
		q := fmt.Sprintf("SELECT * FROM %s LIMIT %d", MyMeasurement, 20)
		res, err = queryDB(clnt, q)
		if err != nil {
			log.Fatal(err)
		}

		for i, row := range res[0].Series[0].Values {
			t, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				log.Fatal(err)
			}
			val := row[1].(string)
			log.Printf("[%2d] %s: %s\n", i, t.Format(time.Stamp), val)
		}
	*/

	//@TODO we should use RethinkDB for this end point
	client := env.GetEnv().Influx
	q := fmt.Sprintf("select * from \"24h\".http_response where ServiceId='%d' order by time desc limit 1 ", 1)
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

}
