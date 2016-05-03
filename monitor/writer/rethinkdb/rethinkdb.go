package rethinkdb

import (
	r "github.com/dancannon/gorethink"

	"fmt"
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

type Writer struct {
	DataChan chan *core.HTTPMetric
	Size     int
	session  *r.Session
	config   *config.Config
}

// NewWriter creates a flusher struct
func NewWriter(config *config.Config) *Writer {
	f := &Writer{
		config: config,
	}
	f.Size = 1000
	f.DataChan = make(chan *core.HTTPMetric, f.Size)
	s, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf("%s:%s", config.RethinkDBHost, config.RethinkDBPort),
		Database: config.RethinkDBName,
		MaxIdle:  10,
		MaxOpen:  10,
		Username: config.RethinkDBUser,
		Password: config.RethinkDBPass,
	})
	if err != nil {
		log.Fatalln("Cannot connec to RethinkDB", err.Error())
	}
	s.SetMaxOpenConns(10)
	f.session = s

	return f
}

func (w *Writer) Name() string {
	return "InfluxDB"
}

func (w *Writer) Write(point *core.HTTPMetric) {
	log.Println("Implement this")
}

func (w *Writer) WriteBatch(points []*core.HTTPMetric) (int, error) {
	bufferPoints := make(buffer, FlushThreshold, FlushThreshold)
	for i, p := range points {
		log.Printf("Got data %v", p.Response.Status)

		bufferPoints[i] = &doc{
			Duration:  float64(p.Response.Duration / time.Millisecond),
			Status:    p.Response.Status,
			Body:      p.Response.Body,
			ServiceId: p.Service.ID,
			Error:     p.Response.Error,
		}
		log.Printf("Add point %v", p.Response.Status)
	}

	res, err := r.DB("notyim").Table("http_response").Insert(bufferPoints).Run(w.session)
	defer res.Close() // Always ensure you close the cursor to ensure connections are not leaked

	if err != nil {
		log.Printf("Fail to flush to RethinkDB %s %v", w.config.InfluxdbHost, err)
		return 0, err
	} else {
		log.Printf("Flush %d points to RethinkDB", len(points))
	}
	return len(points), nil
}
