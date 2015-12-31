package monitor

import (
	"github.com/influxdb/influxdb/client/v2"
	"log"
	"time"
)

type Flusher struct {
	DataChan chan StatusResult
	Size     int
	client   client.Client
}

func NewFlusher(c client.Client) *Flusher {
	f := &Flusher{}
	f.Size = 1000
	f.DataChan = make(chan StatusResult, f.Size)
	if c == nil {
		c, _ = client.NewHTTPClient(client.HTTPConfig{
			Addr:     "http://10.8.0.1:8086",
			Username: "root",
			Password: "root",
		})
	}
	f.client = c

	return f
}

func (f *Flusher) Start() {
	var totalPoint = 0
	var bp client.BatchPoints

	for {
		//@TODO write into influxdb
		if totalPoint == 0 {
			bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  "notyim",
				Precision: "s",
			})
		}

		r := <-f.DataChan

		tags := map[string]string{
			"ServiceId": r.Service.Id,
		}
		fields := map[string]interface{}{
			"Duration": float64(r.Response.Duration / time.Millisecond),
			"Status":   r.Response.Status,
			"Body":     r.Response.Body,
		}
		pt, _ := client.NewPoint("http_response", tags, fields, time.Now())
		bp.AddPoint(pt)
		totalPoint += 1

		if totalPoint >= 500 {
			if err := f.client.Write(bp); err != nil {
				log.Printf("Fail to flush to InfluxDB %v", err)
			} else {
				log.Printf("Flush %d points", totalPoint)
			}
			totalPoint = 0
		}
	}
}
