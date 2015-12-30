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

func NewFlusher() *Flusher {
	f := &Flusher{}
	f.Size = 1000
	f.DataChan = make(chan StatusResult, f.Size)
	f.client, _ = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://192.168.99.100:8086",
		Username: "root",
		Password: "root",
	})

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

		tags := map[string]string{"cpu": "cpu-total"}
		fields := map[string]interface{}{
			"duration": float64(r.Response.Duration / time.Millisecond),
			"status":   r.Response.Status,
			"body":     r.Response.Body,
		}
		pt, _ := client.NewPoint("response", tags, fields, time.Now())
		bp.AddPoint(pt)
		totalPoint += 1

		if totalPoint >= 10 {
			if err := f.client.Write(bp); err != nil {
				log.Printf("Fail to flush to InfluxDB %v", err)
			} else {
				log.Printf("Flush %d points", totalPoint)
			}
			totalPoint = 0
		}
	}
}
