package rethinkdb

import (
	r "github.com/dancannon/gorethink"

	"github.com/notyim/gaia/config"
	"github.com/notyim/gaia/monitor/core"
	"log"
	"time"
)

const (
	// FlushThreshold is the point we need to reach before flushing to storage
	// a smaller value mean more frequently write
	FlushThreshold = 5
)

type doc struct {
	ID        string  `gorethink:"id"`
	Duration  float64 `gorethink:"duration"`
	Status    int     `gorethink:"status"`
	Body      string  `gorethink:"body"`
	ServiceId string  `gorethink:"serviceId"`
	Error     error   `gorethink:"error"`
}

type buffer []*doc

// Flusher represents a flusher that flushs data to a storage backend
type Flusher struct {
	DataChan chan *core.HTTPMetric
	Size     int
	session  *r.Session
	config   *config.Config
	buffer   buffer
}

// NewFlusher creata a flusher struct
func NewFlusher(config *config.Config) *Flusher {
	f := &Flusher{
		config: config,
	}
	f.Size = 1000
	f.DataChan = make(chan *core.HTTPMetric, f.Size)
	s, err := r.Connect(r.ConnectOpts{
		Address:  "locahost:28015",
		Database: "notyim",
		MaxIdle:  10,
		MaxOpen:  10,
	})
	if err != nil {
		log.Fatalln("Cannot connec to RethinkDB", err.Error())
	}
	s.SetMaxOpenConns(10)
	f.session = s

	return f
}

// Start accepts incoming data from its own data channel and flush to backend
func (f *Flusher) Start() {
	var totalPoint = 0

	for {
		if totalPoint == 0 {
			f.buffer = make(buffer, FlushThreshold, FlushThreshold)
		}

		res := <-f.DataChan
		log.Printf("Got data %v", res.Response.Status)

		f.buffer[totalPoint] = &doc{
			Duration:  float64(res.Response.Duration / time.Millisecond),
			Status:    res.Response.Status,
			Body:      res.Response.Body,
			ServiceId: res.Service.ID,
			Error:     res.Response.Error,
		}

		log.Printf("Add point %v", res.Response.Status)
		totalPoint++

		if totalPoint >= FlushThreshold {
			go func(buffer buffer) {
				res, err := r.DB("notyim").Table("http_response").Insert(buffer).Run(f.session)
				if err != nil {
					log.Printf("Fail to flush to RethinkDB %s %v", f.config.InfluxdbHost, err)
				} else {
					log.Printf("Flush %d points", totalPoint)
				}
			}(f.buffer)
			totalPoint = 0
		}
	}
}
