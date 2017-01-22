package models

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/notyim/gaia/db/influxdb"
	"github.com/notyim/gaia/types"
	"log"
	"time"
)

type CheckResult types.HTTPCheckResponse

func (c *CheckResult) Point() *client.Point {
	tags := map[string]string{
		"CheckID": c.CheckID,
	}

	fields := map[string]interface{}{
		"TotalSize": c.TotalSize,
		"TotalTime": c.TotalTime,
		"Error":     c.Error,
	}

	if c.Error {
		fields["ErrorMessage"] = c.ErrorMessage
	}

	point, err := client.NewPoint("http_response", tags, fields, time.Now())

	if err != nil {
		//TODO log error
		log.Println("Cannot create point", err)
		return nil
	}

	return point
}

func (c *CheckResult) Save() {
	log.Println("FLush", c, "to InfluxDB")
	// TODO: get timestamp from client
	if err := influxdb.WritePoint(c.Point()); err != nil {
		log.Println("Cannot write batch points to influxdb", err)
	}
}
