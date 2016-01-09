package influxdb

import (
	"github.com/influxdb/influxdb/client/v2"
	"github.com/notyim/gaia/config"
	"log"
	"time"
)

const (
	// FlushThreshold is the point we need to reach before flushing to storage
	// a smaller value mean more frequently write
	FlushThreshold = 50
)

// Flusher represents a flusher that flushs data to a storage backend
type Flusher struct {
	DataChan chan StatusResult
	Size     int
	client   client.Client
	config   *config.Config
}

// NewFlusher creata a flusher struct
func NewFlusher(config *config.Config, c client.Client) *Flusher {
	f := &Flusher{
		config: config,
	}
	f.Size = 1000
	f.DataChan = make(chan StatusResult, f.Size)
	if c == nil {
		c, _ = client.NewHTTPClient(client.HTTPConfig{
			Addr:     config.InfluxdbHost,
			Username: config.InfluxdbUsername,
			Password: config.InfluxdbPassword,
		})
	}
	f.client = c

	return f
}

// Start accepts incoming data from its own data channel and flush to backend
func (f *Flusher) Start() {
	var totalPoint = 0
	var bp client.BatchPoints

	for {
		//@TODO write into influxdb
		if totalPoint == 0 {
			bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  f.config.InfluxdbDb,
				Precision: "s",
			})
		}

		r := <-f.DataChan

		tags := map[string]string{
			"ServiceId": r.Service.ID,
		}
		fields := map[string]interface{}{
			"Duration": float64(r.Response.Duration / time.Millisecond),
			"Status":   r.Response.Status,
			"Body":     r.Response.Body,
		}

		if nil != r.Response.Error {
			fields["Error"] = r.Response.Error
		}

		pt, _ := client.NewPoint("http_response", tags, fields, time.Now())
		bp.AddPoint(pt)

		//pb, _ := client.NewPoint("http_response_body", tags, fields, time.Now())
		//bp.AddPoint(pt)

		totalPoint++

		if totalPoint >= FlushThreshold {
			if err := f.client.Write(bp); err != nil {
				log.Printf("Fail to flush to InfluxDB %s %v", f.config.InfluxdbHost, err)
			} else {
				log.Printf("Flush %d points", totalPoint)
			}
			totalPoint = 0
		}
	}
}
