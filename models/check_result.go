package models

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/notyim/gaia/db/influxdb"
	"github.com/notyim/gaia/types"
	"log"
	"time"
)

type CheckResult types.HTTPCheckResponse

func (c *CheckResult) TotalTimeAsInt() int {
	return int(c.Time["Total"] / time.Millisecond)
}

func (c *CheckResult) Point() *client.Point {
	tags := map[string]string{
		"CheckID": c.CheckID,
	}

	fields := map[string]interface{}{
		"TotalTime": c.TotalTimeAsInt(),
		"Error":     c.Error,
	}

	if c.Error {
		fields["ErrorMessage"] = c.ErrorMessage
	}

	point, err := client.NewPoint("http_response", tags, fields, c.CheckedAt)

	if err != nil {
		//TODO log error
		log.Println("Cannot create point", err)
		return nil
	}

	return point
}

func (c *CheckResult) Save() {
	log.Println("FLush", c, "to InfluxDB")
	if err := influxdb.WritePoint(c.Point()); err != nil {
		log.Println("Cannot write batch points to influxdb", err)
	}
}
